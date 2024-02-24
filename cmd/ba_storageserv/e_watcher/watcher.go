package watcher

import (
	"context"
	"errors"
	"fmt"
	shard "junodb_lite/pkg/b_shard"
	etcd "junodb_lite/pkg/c_etcd"
	"time"
)

type IWatchEvtHandler interface {
	UpdateShards(shards shard.Map) bool
	UpdateRedistShards(shards shard.Map) bool

	RedistStart(ratelimit int)
	RedistResume(ratelimit int)
	RedistStop()
	RedistIsInProgress() bool
}

type Watcher struct {
	clustername string
	zoneid      uint16
	nodeid      uint16
	version     uint32
	etcdcli     *etcd.Client
	cancel      context.CancelFunc
	updateDelay time.Duration
	name        string
	hdr         IWatchEvtHandler
}

var (
	theWatcher *Watcher
	evtHdr     IWatchEvtHandler
)

func RegisterWatchEvtHandler(hdr IWatchEvtHandler) {
	if theWatcher != nil {
		theWatcher.hdr = hdr
	} else {
		evtHdr = hdr
	}
}

func Init(clustername string, zoneid uint16, nodeid uint16, cfg *etcd.Config, delay time.Duration, version uint32) (err error) {
	etcdcli := etcd.GetEtcdCli()
	if etcdcli == nil {
		etcdcli = etcd.NewEtcdClient(cfg, clustername)
	}

	if etcdcli == nil {
		return errors.New("failed to connect to etcd")
	}

	theWatcher = newWatcher(clustername, zoneid, nodeid, etcdcli, delay, version)
	if evtHdr != nil {
		theWatcher.hdr = evtHdr
	}
	go theWatcher.Watch()
	return nil
}

func newWatcher(clustername string, zoneid uint16, nodeid uint16, cli *etcd.Client, deley time.Duration, version uint32) *Watcher {
	w := &Watcher{
		clustername: clustername,
		zoneid:      zoneid,
		nodeid:      nodeid,
		etcdcli:     cli,
		updateDelay: deley,
		version:     version,
	}

	w.name = fmt.Sprintf("ss-%d-%d watcher", w.zoneid, w.nodeid)
	return w
}

func (w *Watcher) Watch() (cancel context.CancelFunc, err error) {
	return nil, nil
}
