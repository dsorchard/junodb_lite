package etcd

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/prometheus/common/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/namespace"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	errNotInitialized = errors.New("etcd client not initialized")
)

type Client struct {
	initOnce  sync.Once
	config    Config
	keyPrefix string
	client    *clientv3.Client
	doneCh    chan struct{}
	wg        sync.WaitGroup
}

const NotFound = "NotFound"

var (
	shuffleDone = false
	setOnce     sync.Once
)

func NewEtcdClient(cfg *Config, clusterName string) *Client {

	var client *clientv3.Client
	var err error

	now := time.Now()
	m := now.Second() % len(cfg.Endpoints)

	// Shuffle to balance load
	if m > 0 && !shuffleDone {
		endp := make([]string, len(cfg.Endpoints))
		copy(endp[0:], cfg.Endpoints[0:])
		copy(cfg.Endpoints[0:], endp[m:])
		copy(cfg.Endpoints[len(cfg.Endpoints)-m:], endp[0:m])
	}
	shuffleDone = true

	setOnce.Do(func() { // Bypass http_proxy for connecting to etcd server.
		val := strings.Join(cfg.Endpoints, ",")
		curr := os.Getenv("NO_PROXY")
		if strings.Contains(curr, val) {
			return
		}
		if len(curr) > 0 {
			val += "," + curr
		}
		os.Setenv("NO_PROXY", val)
		os.Setenv("no_proxy", val)
	})

	for i := 0; i < cfg.MaxConnectAttempts; i++ {
		client, err = clientv3.New((*cfg).Config)

		if err == nil {
			break
		}

		if client != nil {
			client.Close()
		}

		if i >= cfg.MaxConnectAttempts-1 {
			glog.Warningf("etcd: %v.", err)
			return nil
		}

		glog.Warningf("etcd: %v. Retry ...", err)
		backoff := (i + 1) * 2
		if backoff > cfg.MaxConnectBackoff {
			backoff = cfg.MaxConnectBackoff
		}
		time.Sleep(time.Duration(backoff) * time.Second)
	}

	etcdcli := &Client{
		client: client,
		config: *cfg,
		doneCh: make(chan struct{}),
	}

	etcdcli.keyPrefix = cfg.EtcdKeyPrefix + clusterName + TagCompDelimiter
	etcdcli.client.KV = namespace.NewKV(client.KV, etcdcli.keyPrefix)
	etcdcli.client.Watcher = namespace.NewWatcher(client.Watcher, etcdcli.keyPrefix)
	return etcdcli
}

func (e *Client) Close() {
	close(e.doneCh)

	if e.client != nil {
		e.client.Close()
	}
}

func (e *Client) GetDoneCh() chan struct{} {
	return e.doneCh
}

func (e *Client) GetValue(k string) (value string, err error) {
	var resp *clientv3.GetResponse
	resp, err = e.get(k)
	if err != nil {
		glog.Errorf("%v", err)
		return
	}
	if resp != nil {
		sz := len(resp.Kvs)
		if sz == 1 {
			value = string(resp.Kvs[0].Value)
		} else if sz == 0 {
			err = fmt.Errorf("key '%s' not found.", k)
			value = NotFound
		} else { /// not seem to be possible
			err = fmt.Errorf("unexpected response. %s", k)
		}
	} else {
		err = fmt.Errorf("unexpected nil response. key: %s", k)
	}
	return
}

func (e *Client) get(key string, params ...int) (resp *clientv3.GetResponse, err error) {
	if e.client == nil {
		err = errNotInitialized
		return
	}

	//optional params: maxtries and backoff sleep time
	maxTries := 2
	backoff := 1 * time.Second
	if len(params) > 0 {
		maxTries = params[0]
	}
	if len(params) > 1 {
		backoff = time.Duration(params[1]) * time.Second
	}

	for i := 0; i < maxTries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), e.config.RequestTimeout)
		resp, err = e.client.Get(ctx, key)
		cancel()

		if err == nil {
			return
		}

		glog.Warningf("etcd get: %v. Retry ...", err)
		time.Sleep(backoff)
	}

	if err != nil {
		glog.Errorf("[ERROR]: etcd get: key=%s err=%v", key, err)
	}
	return
}

func (e *Client) PutValue(key string, val string, params ...int) (err error) {
	if e.client == nil {
		err = errNotInitialized
		return
	}

	// optional params: maxtries and backoff sleep time
	maxTries := 1
	backoff := 1 * time.Second // second
	if len(params) > 0 {
		maxTries = params[0]
	}
	if len(params) > 1 {
		backoff = time.Duration(params[1]) * time.Second
	}

	var valStr = val
	var endStr string
	if len(val) >= 50 {
		endStr = " ..."
		valStr = val[:50]
	}

	log.Debugf("etcd put: key=%s%s val=%s%s", e.keyPrefix, key, valStr, endStr)

	for i := 0; i < maxTries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), e.config.RequestTimeout)
		_, err = e.client.Put(ctx, key, val)
		cancel()
		if err == nil {
			break
		}

		if i >= maxTries-1 {
			glog.Errorf("[ERROR]: etcd put: %v", err)
			return err
		}

		glog.Warningf("etcd put: %v. Retry ...", err)
		time.Sleep(backoff)
	}

	return nil
}

// Batch operations of delete and put.
func (e *Client) PutValuesWithTxn(op []clientv3.Op) (err error) {
	if e.client == nil {
		err = errNotInitialized
		return
	}

	if len(op) == 0 {
		return nil
	}

	glog.Infof("etcd txn:")
	for i := 0; i < len(op); i++ {
		if op[i].IsDelete() {
			glog.Infof("etcd delete: beginKey=%s%s", e.keyPrefix, string(op[i].KeyBytes()))
			endKey := op[i].RangeBytes()
			if endKey != nil {
				glog.Infof("               endKey=%s%s", e.keyPrefix, string(endKey))
			}
		} else {
			val := string(op[i].ValueBytes())
			var valStr = val
			var endStr string
			if len(val) >= 20 {
				endStr = " ..."
				valStr = val[:20]
			}
			log.Debugf("etcd put: key=%s%s val=%s%s", e.keyPrefix,
				string(op[i].KeyBytes()), valStr, endStr)
		}
	}
	glog.Infof("ops_count=%d", len(op))

	maxTries := 5
	for i := 0; i < maxTries; i++ {
		timeout := e.config.RequestTimeout
		if timeout < 5*time.Second {
			timeout = 5 * time.Second
		}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		_, err = e.client.Txn(ctx).If().Then(op...).Commit()
		cancel()

		if err == nil {
			break
		}

		if i >= maxTries-1 {
			glog.Errorf("[ERROR]: etcd txn aborted: %v", err)
			return err
		}

		glog.Warningf("etcd: %v. Retry ...", err)
		backoff := 5 * time.Second
		time.Sleep(backoff)
	}

	glog.Infof("etcd txn completed")
	return nil
}

func (e *Client) DeleteKey(key string) (err error) {
	return e.DeleteKeyWithPrefix(key, false)
}

func (e *Client) DeleteKeyWithPrefix(key string, isPrefix ...bool) (err error) {

	usePrefix := true
	if len(isPrefix) > 0 && !isPrefix[0] {
		usePrefix = false
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.config.RequestTimeout)

	glog.Infof("etcd delete: key=%s%s isPrefix=%v", e.keyPrefix, key, usePrefix)
	if usePrefix {
		_, err = e.client.Delete(ctx, key, clientv3.WithPrefix())
	} else {
		_, err = e.client.Delete(ctx, key)
	}
	cancel()

	if err != nil {
		glog.Errorf("%v", err)
	}

	return
}

func (e *Client) WatchEvt(key string, ctx context.Context) (ch clientv3.WatchChan, err error) {
	if e.client == nil {
		err = errNotInitialized
		return
	}

	ch = e.client.Watch(ctx, key, clientv3.WithProgressNotify())
	return
}

func (e *Client) Watch(key string, handler IWatchHandler, opts ...clientv3.OpOption) (cancel context.CancelFunc, err error) {
	if e.client == nil {
		err = errNotInitialized
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	ch := e.client.Watch(ctx, key, opts...)
	e.wg.Add(1)
	go func() {
		glog.Info("start watcher go routine")
		defer e.wg.Done()
		for {
			select {
			case r := <-ch:
				glog.Info("handle event")
				handler.OnEvent(r.Events...)
				//				for i, e := range r.Events {
				//					fmt.Printf("got %d %s", i, e)
				//				}
			case <-ctx.Done():
				glog.Info("Cancel")
				return
			case <-e.doneCh:
				return
			}
		}
	}()
	return cancel, nil
}

// different flavor
func (e *Client) WatchEvents(key string, ch chan int) {
	ctx := context.Background()

	rch := e.client.Watch(ctx, key,
		clientv3.WithProgressNotify())

	glog.Infof("etcd: Watcher waits for events.")
	for {
		select {
		case entry := <-rch:
			for _, ev := range entry.Events {
				val, _ := strconv.Atoi(string(ev.Kv.Value))

				glog.Infof("etcd watch: type=%s key=%q val=%q", ev.Type, ev.Kv.Key, ev.Kv.Value)
				if ev.Type == clientv3.EventTypeDelete {
					continue
				}
				// Notify subscriber
				ch <- val
			}
		case <-ctx.Done():
			glog.Info("Cancel")
			return
		case <-e.doneCh:
			close(ch)
			glog.Infof("etcd: Watcher exits.")
			return
		}
	}
}

// List of transactional operations
type OpList []clientv3.Op
