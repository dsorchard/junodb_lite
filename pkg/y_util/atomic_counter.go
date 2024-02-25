package util

import "sync/atomic"

type AtomicShareCounter struct {
	cnt *uint64
}

type AtomicCounter struct {
	cnt int32
}

func (c *AtomicCounter) Set(cnt int32) {
	atomic.StoreInt32(&c.cnt, cnt)
}
