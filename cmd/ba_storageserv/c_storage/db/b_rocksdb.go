package db

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/tecbot/gorocksdb"
	"io"
	redist "junodb_lite/cmd/ba_storageserv/d_redist"
	shard "junodb_lite/pkg/b_shard"
	"sync/atomic"
)

var rocksdbIndex int32 = 0
var rocksdb [2]IDatabase

func GetDB() IDatabase {
	var index int32 = atomic.LoadInt32(&rocksdbIndex)
	return rocksdb[index]
}

// only called once during start up
func Initialize(
	numShards int, numMicroShards int, numMicroShardGroups int,
	numPrefixDbs int, zoneId int, nodeId int, shardMap shard.Map, lruCacheSizeInMB int) {
	if numMicroShards > 0 {
		//SetEnableMircoShardId(true)
		glog.Infof("Enable micro shards, NumMicroShards=%d, numMshardGroups=%d", numMicroShards, numMicroShardGroups)
	}

	//if DBConfig.NewLRUCacheSizeInMB == 0 && lruCacheSizeInMB > 0 { // Use computed value
	//	DBConfig.NewLRUCacheSizeInMB = lruCacheSizeInMB
	//}
	db := newRocksDB(numShards, numMicroShards, numMicroShardGroups, numPrefixDbs, zoneId, nodeId, shardMap)
	rocksdb[rocksdbIndex] = db
	// safe guard?
	rocksdb[(rocksdbIndex+1)%2] = db
}

type RocksDB struct {
	zoneId       int
	nodeId       int
	numShards    int
	shards       shard.Map // from current shard map
	redistShards shard.Map // from redistribution, temporary/not commited
	sharding     IDBSharding
}

func (r *RocksDB) Setup() {
	//TODO implement me
	panic("implement me")
}

func (r *RocksDB) TruncateExpired() {
	//TODO implement me
	panic("implement me")
}

func (r *RocksDB) Shutdown() {
	//TODO implement me
	panic("implement me")
}

func (r *RocksDB) Put(id RecordID, value []byte) error {
	//TODO implement me
	panic("implement me")
}

func (r *RocksDB) Get(id RecordID, fetchExpired bool) (*Record, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RocksDB) GetRecord(id RecordID, rec *Record) (recExists bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *RocksDB) Delete(id RecordID) error {
	//TODO implement me
	panic("implement me")
}

func (r *RocksDB) IsPresent(id RecordID) (bool, error, *Record) {
	//TODO implement me
	panic("implement me")
}

func (r *RocksDB) IsRecordPresent(id RecordID, rec *Record) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RocksDB) ReplicateSnapshot(shardId shard.ID, rb *redist.Replicator, mshardid int32) bool {
	//TODO implement me
	panic("implement me")
}

func (r *RocksDB) ShardSupported(shardId shard.ID) bool {
	//TODO implement me
	panic("implement me")
}

func (r *RocksDB) UpdateRedistShards(shards shard.Map) {
	//TODO implement me
	panic("implement me")
}

func (r *RocksDB) UpdateShards(shards shard.Map) {
	//TODO implement me
	panic("implement me")
}

func (r *RocksDB) WriteProperty(propKey string, w io.Writer) {
	//TODO implement me
	panic("implement me")
}

func (r *RocksDB) GetIntProperty(propKey string) uint64 {
	//TODO implement me
	panic("implement me")
}

func newRocksDB(numShards int, numMicroShards int, numMicroShardGroups int, numPrefixDbs int, zoneId int, nodeId int, shardMap shard.Map) *RocksDB {
	db := &RocksDB{
		zoneId:    zoneId,
		nodeId:    nodeId,
		numShards: numShards,
		shards:    shardMap,
		sharding:  newDBSharding(numShards, numMicroShards, numMicroShardGroups, numPrefixDbs, fmt.Sprintf("%d-%d", zoneId, nodeId)),
	}
	db.Setup()
	return db
}

// /TODO xuli dbDir...
func newDBSharding(numShards int, numMicroShards int, numMicroShardGroups int, numPrefixDbs int, dbnamePrefix string) (sharding IDBSharding) {
	if numPrefixDbs > 0 { // Use prefix key
		shardFilters := make([]*ShardFilter, numPrefixDbs, numPrefixDbs)
		for i := 0; i < numPrefixDbs; i++ {
			shardFilters[i] = &ShardFilter{shardNum: -1}
		}
		sharding = &ShardingByPrefix{
			dbnamePrefix:        dbnamePrefix,
			dbs:                 make([]*gorocksdb.DB, numPrefixDbs, numPrefixDbs),
			shardFilters:        shardFilters,
			numMicroShards:      numMicroShards,
			numMicroShardGroups: numMicroShardGroups,
		}
	} else { // Not using prefix key
		//sharding = &ShardingByInstance{
		//	dbnamePrefix: dbnamePrefix,
		//	dbs:          make([]*gorocksdb.DB, numShards, numShards),
		//}
	}
	return
}
