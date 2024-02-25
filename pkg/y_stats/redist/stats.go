package redist

import (
	util "junodb_lite/pkg/y_util"
	"strconv"
	"strings"
	"sync/atomic"
)

var (
	StatsTotalCount   = []byte("total")
	StatsOkCount      = []byte("ok")
	StatsErrCount     = []byte("err")
	StatsDropCount    = []byte("drop")
	StatsExpiredCount = []byte("expired")
	StatsMicroShardId = []byte("mshd")
	StatsTagStatus    = []byte("st")
	StatsTagElapse    = []byte("et")

	StatsInProgress    = "P"
	StatsFinish        = "F"
	StatsBegin         = "B"
	StatsAbort         = "A"
	StatsPairDelimiter = "&"
	StatsDelimiter     = "="
)

type KVPairs struct {
	kvpairs map[string]string
}

func NewKVPairs(str string) *KVPairs {
	p := &KVPairs{
		kvpairs: make(map[string]string),
	}
	pairs := strings.Split(str, StatsPairDelimiter)
	for _, kv := range pairs {
		v := strings.Split(kv, StatsDelimiter)
		if len(v) >= 2 {
			p.kvpairs[v[0]] = v[1]
		}
	}
	return p
}

func (p *KVPairs) GetValue(key string, defv string) string {
	value := p.kvpairs[key]
	if len(value) == 0 {
		return defv
	}
	return value
}

func (p *KVPairs) GetInt(key string, defv int) int {
	value := p.kvpairs[key]
	if len(value) == 0 {
		return defv
	}
	s, err := strconv.Atoi(value)
	if err != nil {
		return defv
	}
	return s
}

type Stats struct {
	totalCnt util.AtomicCounter
	okCnt    util.AtomicCounter
	failCnt  util.AtomicCounter
	dropCnt  util.AtomicCounter
	expCnt   util.AtomicCounter
	mshardId int32
	status   string

	// for tracking last successful checkpoint
	lastTotalCnt int32
	lastOkCnt    int32
	lastFailCnt  int32
	lastDropCnt  int32
	lastExpCnt   int32
	lastMShardId int32
}

func NewStats(str string) *Stats {
	kvs := NewKVPairs(str)
	st := &Stats{}
	st.status = kvs.GetValue(string(StatsTagStatus), "")
	st.SetMShardId(int32(kvs.GetInt(string(StatsMicroShardId), 0)))
	st.totalCnt.Set(int32(kvs.GetInt(string(StatsTotalCount), 0)))
	st.okCnt.Set(int32(kvs.GetInt(string(StatsOkCount), 0)))
	st.failCnt.Set(int32(kvs.GetInt(string(StatsErrCount), 0)))
	st.dropCnt.Set(int32(kvs.GetInt(string(StatsDropCount), 0)))
	st.expCnt.Set(int32(kvs.GetInt(string(StatsExpiredCount), 0)))

	//st.SaveCheckPoint()
	return st
}
func (r *Stats) SetMShardId(id int32) {
	r.lastMShardId = r.mshardId
	atomic.StoreInt32(&r.mshardId, id)
}

func (r *Stats) GetStatus() string {
	return r.status
}

func (r *Stats) GetMShardId() int32 {
	return atomic.LoadInt32(&r.lastMShardId)
}

func (r *Stats) Restore(ch *Stats) {
	r.lastTotalCnt = ch.lastTotalCnt
	r.lastOkCnt = ch.lastOkCnt
	r.lastFailCnt = ch.lastFailCnt
	r.lastDropCnt = ch.lastDropCnt
	r.lastExpCnt = ch.lastExpCnt
	r.lastMShardId = ch.lastMShardId
	//r.RestoreFromCheckPoint()
}
