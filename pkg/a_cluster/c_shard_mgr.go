package cluster

import (
	"junodb_lite/pkg/io"
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
