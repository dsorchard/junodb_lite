package db

import (
	gorocksdb "github.com/tecbot/gorocksdb"
	"io"
	redist "junodb_lite/cmd/ba_storageserv/d_redist"
	shard "junodb_lite/pkg/b_shard"
)

type IDBSharding interface {
	getDbInstanceAndKey(id RecordID) (dbInst *gorocksdb.DB, dbKey []byte)

	setupShards(dbnamePrefix string, shardMap shard.Map)

	shutdownShards([]shard.ID)
	shutdown()

	writeProperty(propKey string, w io.Writer)
	getIntProperty(propKey string) uint64

	decodeStorageKey(sskey []byte) ([]byte, []byte, error)
	duplicate() IDBSharding

	replicateSnapshot(shardId shard.ID, rb *redist.Replicator, mshardid int32) bool
}
type ShardingBase struct {
}

func (s *ShardingBase) waitForFinish(rb *redist.Replicator) bool {
	return true
}
