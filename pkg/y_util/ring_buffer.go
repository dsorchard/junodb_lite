package util

import "time"

type QueItem interface {
	OnCleanup()
	OnExpiration()
	Deadline() (deadline time.Time)
	ResetDeadline()
	SetId(id uint32)
	GetId() uint32
	SetInUse(flag bool)
	SetQueTimeout(t time.Duration)
	GetQueTimeout() (t time.Duration)
	IsInUse() bool
}
type RingBuffer struct {
	head     uint32 // Atomic access, updated by reader, used
	tail     uint32 // Atomic access, updated by writer, unused
	capacity uint32 // qsize + 10% extra + 1
	buf      []QueItem
	seqId    uint32
	qsize    uint32 // qsize exposed to user
	extra    uint32
	cursize  int32
}

// drain everything in the ringbuffer
func (rb *RingBuffer) CleanAll() {

}

type QueItemBase struct {
	id       uint32
	flag     uint32
	timeout  time.Duration
	deadline time.Time
}
