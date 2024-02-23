package cluster

import (
	shard "junodb_lite/pkg/b_shard"
	"junodb_lite/pkg/y_conn_mgr"
	util "junodb_lite/pkg/y_util"
	"sync/atomic"
)

var (
	ClusterInfo   [2]Cluster
	shardMgrPair  [2]*ShardManager
	shardMgrIndex int32 = 0
)

type (
	ZoneMarkDown struct {
		markdownid int32
	}
	OutboundSSProcessor struct {
		io.OutboundProcessor
		zoneId      int
		indexInZone int
	}
	ShardManager struct {
		AlgVersion uint32 // default
		shardMap   ShardMap
		connInfo   [][]string
		ssconfig   *io.OutboundConfig
		processors [][]*OutboundSSProcessor
	}
)

func Initialize(args ...interface{}) (err error) {
	return nil
}

func Finalize() {
	GetShardMgr().Shutdown(nil)
}

// Return shardmgr that is active.
func GetShardMgr() *ShardManager {
	var ix int32 = atomic.LoadInt32(&shardMgrIndex)
	return shardMgrPair[ix]
}

func newShardManager(ccfg *Cluster, conf *io.OutboundConfig, statscfg *StatsConfig, curMgr *ShardManager) (m *ShardManager, err error) {
	return nil, nil
}

func (p *ShardManager) newAndStartSSProcessor(zoneId int, indexInZone int, enableBounce bool) *OutboundSSProcessor {
	proc := &OutboundSSProcessor{zoneId: zoneId, indexInZone: indexInZone}
	proc.Init(io.ServiceEndpoint{Addr: p.connInfo[zoneId][indexInZone]}, p.ssconfig, enableBounce)
	proc.SetConnEventHandler(proc)
	proc.Start()
	return proc
}

// Used by admin worker
func (p *ShardManager) GetProcessorsByKey(key []byte) (shardId shard.ID, procs []*OutboundSSProcessor, err error) {
	shardId = shard.ID(util.GetPartitionId(key, uint32(p.shardMap.cluster.NumShards)))
	procs, err = p.GetProcessors(shardId)
	return
}

// Used by admin worker
func (p *ShardManager) GetProcessors(partId shard.ID) ([]*OutboundSSProcessor, error) {

	// for testing connectivity, order does not matter, so always use 1 for now
	zones, nodes, err := p.shardMap.GetNodes(uint32(partId), 1)
	if err != nil {
		return nil, err
	}

	procs := make([]*OutboundSSProcessor, len(zones))
	for i := range zones {
		s := p.processors[int(zones[i])][nodes[i]]
		procs[i] = s
	}

	return procs, nil
}

func (p *ShardManager) Shutdown(curMgr *ShardManager) {

}
