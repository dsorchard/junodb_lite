package io

import (
	"bufio"
	util "junodb_lite/pkg/y_util"
	"net"
	"sync"
)

type OutboundConnector struct {
	id     int
	conn   net.Conn // tcp connection to downstream server
	reader *bufio.Reader

	//reqCh      chan IRequestContext
	//reqPending *util.RingBuffer // lock-free ring buffer

	doneCh    chan struct{}
	monitorCh chan int
	wg        sync.WaitGroup
	closeOnce sync.Once
	config    *OutboundConfig
	state     int32
	//hshaker     IHandshaker
	displayName string
	reqPending  *util.RingBuffer // lock-free ring buffer
	reqCh       chan IRequestContext
}

// wait for all go routine to finish
func (p *OutboundConnector) Shutdown() {
	p.Close()
	p.wg.Wait()
	p.cleanPending()
}
func (p *OutboundConnector) Close() {
	p.closeOnce.Do(func() {
		close(p.doneCh)
		//p.conn.Close()
		p.monitorCh <- p.id
	})
}
func (p *OutboundConnector) cleanPending() {
	// drain the ring buffer
	p.reqPending.CleanAll()
}

// constructor/factory
// each OutboundConnector object will have two go routines
// one for Read; one for Write
func NewOutboundConnector(id int, c net.Conn, reqCh chan IRequestContext, monCh chan int,
	config *OutboundConfig) (p *OutboundConnector) {
	p = &OutboundConnector{
		id:   id,
		conn: c,
		//reader: util.NewBufioReader(c, config.IOBufSize),
		reqCh: reqCh,
		//reqPending: util.NewRingBufferWithExtra(uint32(config.MaxPendingQueSize-1),
		//	uint32(config.PendingQueExtra)),
		doneCh:    make(chan struct{}),
		monitorCh: monCh,
		config:    config,
	}

	return p
}
