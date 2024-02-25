package redist

import (
	"errors"
	"github.com/golang/glog"
	proto "junodb_lite/pkg/ac_proto"
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

func (r *Replicator) RestoreSnapShotState(s *redistst.Stats) {
	r.snapshotStats.Restore(s)
}

func (r *Replicator) SendRequest(msg *proto.RawMessage, params ...bool) error {
	glog.Infof("redist:ReplicateRequest: proc=%v, rb=%v", r.processor, r)
	if r.processor == nil {
		return errors.New("outbound processor is not available")
	}

	//default
	realtime := true
	cntOnFailure := true
	if len(params) > 0 {
		realtime = params[0]
	}

	if len(params) > 1 {
		cntOnFailure = params[1]
	}

	var stats *redistst.Stats = &r.snapshotStats
	if realtime {
		stats = &r.realtimeStats
	}

	reqctx := NewRedistRequestContext(msg, r.processor.GetRequestCh(), stats)
	var err error

	if realtime {
		err = r.processor.SendRequest(reqctx)
	} else {
		//err = r.processor.SendRequestLowPriority(reqctx)
	}

	if err == nil {
		//stats.IncreaseTotalCnt()
		return nil
	}

	// forwarding queue is full or not ready
	if cntOnFailure {
		//stats.IncreaseTotalCnt()
		//stats.IncreaseDropCnt()
	}
	return errors.New("Forwarding queue is either full or not ready, drop req")
}
