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
