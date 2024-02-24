package redist

import (
	shard "junodb_lite/pkg/b_shard"
	etcd "junodb_lite/pkg/c_etcd"
	io "junodb_lite/pkg/y_conn_mgr"
	"sync"
)

type Replicator struct {
	shardId   shard.ID
	processor *io.OutboundProcessor
	wg        *sync.WaitGroup
	//snapshotStats redistst.Stats
	//realtimeStats redistst.Stats
	statskey  string
	ratelimit int
	etcdcli   *etcd.Client
}
