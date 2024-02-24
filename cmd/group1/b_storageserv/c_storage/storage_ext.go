package storage

import (
	"junodb_lite/cmd/group1/b_storageserv/c_storage/db"
	redist "junodb_lite/cmd/group1/b_storageserv/d_redist"
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
	//TODO implement me
	panic("implement me")
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
