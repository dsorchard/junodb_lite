package db

import proto "junodb_lite/pkg/ac_proto"

type (
	compactionFilter struct {
		shardFilter *ShardFilter
	}
	recordFlagT byte

	valueHolderI interface {
		Data() []byte
		Size() int
		Free()
	}

	RecordHeader struct {
		Version              uint32
		CreationTime         uint32
		LastModificationTime uint64
		ExpirationTime       uint32

		OriginatorRequestId proto.RequestId
		RequestId           proto.RequestId
		flag                recordFlagT
	}
	Record struct {
		RecordHeader
		Payload proto.Payload
		holder  valueHolderI
	}
)
