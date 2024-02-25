package redist

import (
	"errors"
	"github.com/golang/glog"
	shard "junodb_lite/pkg/b_shard"
	etcd "junodb_lite/pkg/c_etcd"
	"sync"
	"sync/atomic"
)

type IDBRedistHandler interface {
	SendRedistSnapshot(shardId shard.ID, rb *Replicator, mshardid int32) bool
}

var (
	enabled    int32
	theManager *Manager
	theLock    sync.RWMutex // or use unsafe ponter?
	redistHdr  IDBRedistHandler
	etcdcli    *etcd.Client
)

func RegisterDBRedistHandler(hdr IDBRedistHandler) {
	redistHdr = hdr
}

func Init(clustername string, zoneid uint16, nodeid uint16, cfg *etcd.Config) (err error) {
	glog.Infof("redist.Init: zoneid:%d, nodeid:%d", zoneid, nodeid)

	etcdcli := etcd.GetEtcdCli()
	if etcdcli == nil {
		etcdcli = etcd.NewEtcdClient(cfg, clustername)
	}

	if etcdcli == nil {
		return errors.New("failed to connect to etcd")
	}
	return nil
}

func SetManager(mgr *Manager) {
	theLock.Lock()
	theManager = mgr
	defer theLock.Unlock()
}
func SetEnable(flag bool) {
	glog.Infof("RedistEnabled: %v", flag)
	var value int32 = 0
	if flag {
		value = 1
	}
	atomic.StoreInt32(&enabled, value)
}
