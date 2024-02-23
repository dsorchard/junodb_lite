package util

import (
	"encoding/binary"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/spaolacci/murmur3"
	"time"
)

const uuidEpoch = 122192928000000000 // UUID epoch (October 15, 1582)

func Murmur3Hash(data []byte) uint32 {
	return murmur3.Sum32(data)
}

func GetPartitionId(key []byte, numShards uint32) uint16 {
	hashcode := Murmur3Hash(key)
	return uint16(hashcode % numShards)
}

func GetTimeFromUUIDv1(id uuid.UUID) (tm time.Time, err error) {
	if id[6]&0xF0 != 0x10 {
		err = fmt.Errorf("not v1 UUID")
		return
	}
	var buf [8]byte
	buf[0] = id[6] & 0xF
	buf[1] = id[7]
	buf[2] = id[4]
	buf[3] = id[5]
	copy(buf[4:], id[:4])

	timestamp := (binary.BigEndian.Uint64(buf[:]) - uuidEpoch) * 100
	tm = time.Unix(0, int64(timestamp))
	return
}
