package cluster

import (
	"junodb_lite/pkg/z_conn_mgr"
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
		z_conn_mgr.OutboundProcessor
		zoneId      int
		indexInZone int
	}
	ShardManager struct {
		AlgVersion uint32 // default
		shardMap   ShardMap
		connInfo   [][]string
		ssconfig   *z_conn_mgr.OutboundConfig
		processors [][]*OutboundSSProcessor
	}
)

func Initialize(args ...interface{}) (err error) {
	return nil
}

func newShardManager(ccfg *Cluster, conf *z_conn_mgr.OutboundConfig, statscfg *StatsConfig, curMgr *ShardManager) (m *ShardManager, err error) {
	return nil, nil
}

func (p *ShardManager) newAndStartSSProcessor(zoneId int, indexInZone int, enableBounce bool) *OutboundSSProcessor {
	proc := &OutboundSSProcessor{zoneId: zoneId, indexInZone: indexInZone}
	proc.Init(z_conn_mgr.ServiceEndpoint{Addr: p.connInfo[zoneId][indexInZone]}, p.ssconfig, enableBounce)
	proc.SetConnEventHandler(proc)
	proc.Start()
	return proc
}
