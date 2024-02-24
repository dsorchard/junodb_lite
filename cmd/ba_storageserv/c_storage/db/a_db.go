package db

import (
	"io"
	redist "junodb_lite/cmd/ba_storageserv/d_redist"
	shard "junodb_lite/pkg/b_shard"
)

type IDatabase interface {
	Setup()
	TruncateExpired()
	Shutdown()

	Put(id RecordID, value []byte) error
	Get(id RecordID, fetchExpired bool) (*Record, error)
	GetRecord(id RecordID, rec *Record) (recExists bool, err error)
	Delete(id RecordID) error

	IsPresent(id RecordID) (bool, error, *Record)
	IsRecordPresent(id RecordID, rec *Record) (bool, error)

	ReplicateSnapshot(shardId shard.ID, r *redist.Replicator, mshardid int32) bool
	ShardSupported(shardId shard.ID) bool
	UpdateRedistShards(shards shard.Map)
	UpdateShards(shards shard.Map)

	WriteProperty(propKey string, w io.Writer)
	GetIntProperty(propKey string) uint64
}
