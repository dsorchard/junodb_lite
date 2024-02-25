package redist

import (
	"errors"
	"github.com/golang/glog"
	shard "junodb_lite/pkg/b_shard"
	etcd "junodb_lite/pkg/c_etcd"
	io "junodb_lite/pkg/y_conn_mgr"
	redistst "junodb_lite/pkg/y_stats/redist"
	"sync"
	"sync/atomic"
	"time"
)

type Manager struct {
	zoneid       uint16
	nodeid       uint16
	ssconfig     *io.OutboundConfig
	nodeConnInfo []string                 // connection info for the all nodes in the rack
	processors   []*io.OutboundProcessor  // processors corresponding to new nodes
	changeMap    map[shard.ID]*Replicator // shards need to move to new node
	etcdcli      *etcd.Client
	wg           sync.WaitGroup
	stop         int32 // atomic flag to signal Manager to stop: 1 - stop, 0 - ok
	redistDone   int32 // atomic flag to indicate redistribution is done (all snapshot transferred)
}

func NewManager(zoneid uint16, nodeid uint16, connInfo []string,
	changeMap map[uint16]uint16, conf *Config, cli *etcd.Client, ratelimit int) (m *Manager, err error) {

	glog.Infof("NewManager: %v, %v", connInfo, changeMap)
	m = &Manager{
		zoneid:       zoneid,
		nodeid:       nodeid,
		nodeConnInfo: connInfo,
		ssconfig:     &conf.Outbound,
		processors:   make([]*io.OutboundProcessor, len(connInfo)),
		changeMap:    make(map[shard.ID]*Replicator),
		etcdcli:      cli,
		stop:         0,
		redistDone:   0,
	}

	for shardid, nid := range changeMap {
		if nid >= uint16(len(connInfo)) {
			return nil, errors.New("bad Shard Change Map: node id out of bound")
		}

		processor := m.processors[nid]

		if processor == nil {
			processor = io.NewOutboundProcessor(connInfo[nid], &conf.Outbound, false)
			if processor == nil {
				return nil, errors.New("bad redistr manager: failed to create processor")
			}
			m.processors[nid] = processor
		}
		glog.Infof("processor created: %v", processor)
		statskey := etcd.KeyRedistNodeState(int(m.zoneid), int(m.nodeid), int(shardid))
		Replicator := NewBalancer(shard.ID(shardid), processor, &m.wg, statskey, ratelimit, cli)
		m.changeMap[shard.ID(shardid)] = Replicator
	}

	glog.Infof("rebalance mananger: %v", m)
	m.wg.Add(1)
	go m.Start()
	return m, nil
}

func (m *Manager) Start() {
	defer m.wg.Done()
	defer atomic.StoreInt32(&m.redistDone, 1)
	// wait a little bit so that the outboundconnectors are ready to use.
	time.Sleep(1 * time.Second)
	// TODO manage concurrent Replicator snapshot forwarding & rate_limit
	totalShards := len(m.changeMap)

	// try 5 times
	for i := 0; i < 5; i++ {
		finishCnt := 0

		for _, rb := range m.changeMap {
			//m.wg.Add(1)
			//go snapshotHdr.Send(rb.GetShardId(), rb, &m.wg)
			// for now, no parallization on redistributing the snapshot

			// TO revisit
			//if m.IsStopped() {
			//	break
			//}

			// check for snapshot finsih
			key := etcd.KeyRedistNodeState(int(m.zoneid), int(m.nodeid), int(rb.GetShardId()))
			curval, err := m.etcdcli.GetValue(key)
			var mshardid int32 = 0
			if err == nil {
				st := redistst.NewStats(curval)

				status := st.GetStatus()
				if status == redistst.StatsFinish {
					// already finished, skip
					glog.Infof("%d completed, skip", int(rb.GetShardId()))
					finishCnt++
					continue
				} else if status == redistst.StatsAbort {
					// resume from next mshard id
					rb.RestoretSnapShotState(st)
					mshardid = st.GetMShardId()
					if mshardid != 0 {
						mshardid++
					}
				}
			}

			// send the shard
			redistHdr.SendRedistSnapshot(rb.GetShardId(), rb, mshardid)

			curval, err = m.etcdcli.GetValue(key)
			if err == nil {
				st := redistst.NewStats(curval)
				if st.GetStatus() == redistst.StatsFinish {
					// already finished, skip
					finishCnt++
					continue
				}
			}
		}

		if finishCnt == totalShards {
			glog.Infof("Redistribution finished: total %d shards", totalShards)
			return
		}
	}
	glog.Infof("Redistribution aborted -- too many errors")
}
