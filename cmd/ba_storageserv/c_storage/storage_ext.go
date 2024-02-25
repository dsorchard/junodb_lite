package storage

import (
	"github.com/golang/glog"
	"junodb_lite/cmd/ba_storageserv/c_storage/db"
	redist "junodb_lite/cmd/ba_storageserv/d_redist"
	cluster "junodb_lite/pkg/b_cluster"
	shard "junodb_lite/pkg/b_shard"
	etcd "junodb_lite/pkg/c_etcd"
)

type SSRedistWatchHandler struct {
	zoneid  uint16
	nodeid  uint16
	etcdcli *etcd.Client
}

func (h *SSRedistWatchHandler) SendRedistSnapshot(shardId shard.ID, rb *redist.Replicator, mshardid int32) bool {
	//TODO implement me
	panic("implement me")
}

func (h *SSRedistWatchHandler) RedistStart(ratelimit int) {

	// get clusterinfo
	rw := etcd.GetClsReadWriter()
	if rw == nil {
		glog.Info("no etcd reader")
		return
	}

	// clster contains new node info
	var clster cluster.Cluster
	_, err := clster.ReadWithRedistInfo(rw)
	if err != nil {
		glog.Infof("nodeinfo, %s", err.Error())
		return
	}

	changeMap, err := rw.ReadRedistChangeMap(int(h.zoneid), int(h.nodeid))
	if err != nil {
		// if changemap is empty, no need to do anything
		glog.Infof("changemap, %s", err.Error())
		return
	}

	glog.Infof("redist process start")
	glog.Infof("redist nodeinfo: %#v, change map: %v", clster.ConnInfo[h.zoneid], changeMap)

	mgr, err := redist.NewManager(h.zoneid, h.nodeid, clster.ConnInfo[h.zoneid],
		changeMap, &redist.RedistConfig, h.etcdcli, ratelimit)

	if err != nil {
		return
	}

	// enable the new node
	//
	redist.SetManager(mgr)
	redist.SetEnable(true)
}
func (h *SSRedistWatchHandler) RedistResume(ratelimit int) {
	//TODO implement me
	panic("implement me")
}

func (h *SSRedistWatchHandler) RedistStop() {
	//TODO implement me
	panic("implement me")
}

func (h *SSRedistWatchHandler) RedistIsInProgress() bool {
	//TODO implement me
	panic("implement me")
}

func (h *SSRedistWatchHandler) UpdateShards(shards shard.Map) bool {
	db.GetDB().UpdateShards(shards)
	return true
}

func (h *SSRedistWatchHandler) UpdateRedistShards(shards shard.Map) bool {
	db.GetDB().UpdateRedistShards(shards)
	return true
}
