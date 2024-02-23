package z_conn_mgr

import (
	"bufio"
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
}
