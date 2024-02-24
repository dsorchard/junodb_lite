package storage

import (
	"fmt"
	"github.com/golang/glog"
	config "junodb_lite/cmd/group1/b_storageserv/b_config"
	"junodb_lite/cmd/group1/b_storageserv/c_storage/db"
	redist "junodb_lite/cmd/group1/b_storageserv/d_redist"
	watcher "junodb_lite/cmd/group1/b_storageserv/e_watcher"
	shard "junodb_lite/pkg/b_shard"
	etcd "junodb_lite/pkg/c_etcd"
	"sync"
)

var (
	shardMap shard.Map ///TODO to be removed. it does not seem to be used
	onceSe   sync.Once
)

func Initialize(args ...interface{}) (err error) {
	sz := len(args)
	if sz < 3 {
		err = fmt.Errorf("three arguments expected")
		glog.Error(err)
		return
	}
	var zoneId, machineId, lruCacheSizeInMB int
	var ok bool
	if zoneId, ok = args[0].(int); !ok {
		err = fmt.Errorf("zoneId of type int expected")
		glog.Error(err)
		return
	}
	if machineId, ok = args[1].(int); !ok {
		err = fmt.Errorf("machineId of type int expected")
		glog.Error(err)
		return
	}
	if lruCacheSizeInMB, ok = args[2].(int); !ok {
		err = fmt.Errorf("lruCacheSizeInMB of type int expected")
		glog.Error(err)
		return
	}
	initialize(zoneId, machineId, lruCacheSizeInMB)
	return
}

func Finalize() {
	Shutdown()
}

func initialize(zoneId int, machineId int, lruCacheSizeInMB int) {
	//glog.Debugf("setting up storage engine ...")
	onceSe.Do(func() {
		glog.Info("creating a storage engine instance.")
		shardMap = config.ServerConfig().NewShardMap(zoneId, machineId) ///TODO ...
		cfg := config.ServerConfig()

		InitializeCMap(int(cfg.ClusterInfo.NumShards))
		db.Initialize(int(cfg.ClusterInfo.NumShards), int(cfg.NumMicroShards),
			int(cfg.NumMicroShardGroups), int(cfg.NumPrefixDbs),
			zoneId, machineId, shardMap, lruCacheSizeInMB)

		glog.Infof("storage engine initialized")
	})

	etcdcli := etcd.GetEtcdCli()
	//	if etcdcli == nil {
	//		etcdcli = etcd.NewEtcdClient(cfg, clustername)
	//	}
	handler := SSRedistWatchHandler{zoneid: uint16(zoneId), nodeid: uint16(machineId), etcdcli: etcdcli}
	watcher.RegisterWatchEvtHandler(&handler)
	redist.RegisterDBRedistHandler(&handler)
}
func Shutdown() {
	glog.Infof("shutting down storage engine ...")
	db.GetDB().Shutdown()
}
