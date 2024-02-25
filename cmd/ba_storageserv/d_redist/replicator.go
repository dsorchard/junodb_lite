package redist

import (
	shard "junodb_lite/pkg/b_shard"
	etcd "junodb_lite/pkg/c_etcd"
	io "junodb_lite/pkg/y_conn_mgr"
	redistst "junodb_lite/pkg/y_stats/redist"
	"sync"
)

type Replicator struct {
	shardId       shard.ID
	processor     *io.OutboundProcessor
	wg            *sync.WaitGroup
	snapshotStats redistst.Stats
	realtimeStats redistst.Stats
	statskey      string
	ratelimit     int
	etcdcli       *etcd.Client
}

func NewBalancer(shardId shard.ID, processor *io.OutboundProcessor, wg *sync.WaitGroup, key string, ratelimit int, cli *etcd.Client) (r *Replicator) {
	r = &Replicator{
		shardId:   shardId,
		processor: processor,
		wg:        wg,
		statskey:  key,
		ratelimit: ratelimit,
		etcdcli:   cli,
	}

	//glog.Infof("create new balancer %p, %v", r, r)
	return r
}

func (r *Replicator) GetShardId() shard.ID {
	return r.shardId
}

func (r *Replicator) RestoretSnapShotState(s *redistst.Stats) {
	r.snapshotStats.Restore(s)
}
