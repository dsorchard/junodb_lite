package db

import (
	"github.com/tecbot/gorocksdb"
	"io"
	redist "junodb_lite/cmd/group1/b_storageserv/d_redist"
	shard "junodb_lite/pkg/b_shard"
)

type ShardFilter struct {
	shardNum int32
}

type ShardingByPrefix struct {
	ShardingBase
	DbNames     []string
	PrefixBytes int

	dbnamePrefix        string
	dbs                 []*gorocksdb.DB
	shardFilters        []*ShardFilter // For ComactRangeByShard
	numMicroShards      int
	numMicroShardGroups int
}

func (s *ShardingByPrefix) getDbInstanceAndKey(id RecordID) (dbInst *gorocksdb.DB, dbKey []byte) {
	//TODO implement me
	panic("implement me")
}

func (s *ShardingByPrefix) setupShards(dbnamePrefix string, shardMap shard.Map) {
	//TODO implement me
	panic("implement me")
}

func (s *ShardingByPrefix) shutdownShards(ids []shard.ID) {
	//TODO implement me
	panic("implement me")
}

func (s *ShardingByPrefix) shutdown() {
	//TODO implement me
	panic("implement me")
}

func (s *ShardingByPrefix) writeProperty(propKey string, w io.Writer) {
	//TODO implement me
	panic("implement me")
}

func (s *ShardingByPrefix) getIntProperty(propKey string) uint64 {
	//TODO implement me
	panic("implement me")
}

func (s *ShardingByPrefix) decodeStorageKey(sskey []byte) ([]byte, []byte, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ShardingByPrefix) duplicate() IDBSharding {
	//TODO implement me
	panic("implement me")
}

func (s *ShardingByPrefix) replicateSnapshot(shardId shard.ID, rb *redist.Replicator, mshardid int32) bool {
	//TODO implement me
	panic("implement me")
}
