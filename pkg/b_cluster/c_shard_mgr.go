package cluster

import (
	shard "junodb_lite/pkg/b_shard"
	"junodb_lite/pkg/y_conn_mgr"
	util "junodb_lite/pkg/y_util"
)

var (
	ClusterInfo  [2]Cluster
	shardMgrPair [2]*ShardManager
)

type (
	ZoneMarkDown struct {
		markdownid int32
	}
	OutboundSSProcessor struct {
		y_conn_mgr.OutboundProcessor
		zoneId      int
		indexInZone int
	}
	ShardManager struct {
		AlgVersion uint32 // default
		shardMap   ShardMap
		connInfo   [][]string
		ssconfig   *y_conn_mgr.OutboundConfig
		processors [][]*OutboundSSProcessor
	}
)

func Initialize(args ...interface{}) (err error) {
	return nil
}

func newShardManager(ccfg *Cluster, conf *y_conn_mgr.OutboundConfig, statscfg *StatsConfig, curMgr *ShardManager) (m *ShardManager, err error) {
	return nil, nil
}

func (p *ShardManager) newAndStartSSProcessor(zoneId int, indexInZone int, enableBounce bool) *OutboundSSProcessor {
	proc := &OutboundSSProcessor{zoneId: zoneId, indexInZone: indexInZone}
	proc.Init(y_conn_mgr.ServiceEndpoint{Addr: p.connInfo[zoneId][indexInZone]}, p.ssconfig, enableBounce)
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
