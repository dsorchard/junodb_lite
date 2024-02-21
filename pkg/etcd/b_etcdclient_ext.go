package etcd

import (
	"context"
	"github.com/golang/glog"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strconv"
	"time"
)

func (e *Client) GetVersion() (version, algver uint32, err error) {
	var str string

	// get version
	str, err = e.GetValue(TagVersion)
	if err != nil {
		glog.Warningf("%v", err)
		return
	}

	var n uint64
	n, err = strconv.ParseUint(str, 10, 32)
	if err != nil {
		return
	}
	version = uint32(n)

	// get alg version, default is 1
	algver = 1
	str, err = e.GetValue(TagAlgVer)
	if err != nil {
		if str == NotFound {
			err = nil
		}
		return
	}

	n, err = strconv.ParseUint(str, 10, 32)
	if err != nil {
		return
	}
	algver = uint32(n)

	return
}

func (e *Client) GetUint32(k string) (value uint32, err error) {
	var str string
	str, err = e.GetValue(k)
	if err != nil {
		glog.Warningf("%v", err)
		return
	}
	var n uint64
	n, err = strconv.ParseUint(str, 10, 32)
	if err != nil {
		return
	}
	value = uint32(n)
	return
}

// key is sorted in descending order
func (e *Client) getWithPrefix(key string, params ...int) (resp *clientv3.GetResponse, err error) {
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
		resp, err = e.client.KV.Get(ctx, key, clientv3.WithPrefix(),
			clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
		cancel()

		if err == nil {
			return
		}

		glog.Warningf("etcd get: %v. Retry ...", err)
		time.Sleep(backoff)
	}

	if err != nil {
		glog.Errorf("[ERROR]: etcd get: key=%s%s err=%v", e.keyPrefix, key, err)
	}
	return
}
