package redist

import (
	"errors"
	"github.com/golang/glog"
	shard "junodb_lite/pkg/b_shard"
	etcd "junodb_lite/pkg/c_etcd"
	io "junodb_lite/pkg/y_conn_mgr"
	"sync"
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

type Manager struct {
	zoneid       uint16
	nodeid       uint16
	ssconfig     *io.OutboundConfig
	nodeConnInfo []string                 // connection info for the all nodes in the rack
	processors   []*io.OutboundProcessor  // processors corresponding to new nodes
	changeMap    map[shard.ID]*Replicator // shards need to move to new node
	etcdcli      *etcd.Client
	wg           sync.WaitGroup
	stop         int32 // atomic flag to signal Manager to stop: 1 - stop, 0 - ok
	redistDone   int32 // atomic flag to indicate redistribution is done (all snapshot transferred)
}

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
