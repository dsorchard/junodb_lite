package storage

import "sync"

var (
	///TODO to use sharded sync.Map
	prepareMap []*sync.Map // sharded map
)

func InitializeCMap(numShards int) {
	prepareMap = make([]*sync.Map, numShards)
	for i := 0; i < numShards; i++ {
		prepareMap[i] = new(sync.Map) // it's a bit wast for now.
	}
}
