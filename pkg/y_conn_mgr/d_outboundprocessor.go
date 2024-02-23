package y_conn_mgr

import (
	"sync"
	"time"
)

type (
	IConnEventHandler interface {
		OnConnectSuccess(conn Conn, connector *OutboundConnector, timeTaken time.Duration)
		OnConnectError(timeTaken time.Duration, connStr string, err error)
	}

	// OutboundProcessor manages a pool of one or more underlying connections
	// to a downstream server; It also bounces incoming requests when all
	// connections are down.
	//
	OutboundProcessor struct {
		wg         sync.WaitGroup
		connEvHdlr IConnEventHandler
	}
)

func (p *OutboundProcessor) Start() {
	p.wg.Add(1)
	go p.Run()
}
func (p *OutboundProcessor) Run() {
	defer p.wg.Done()
}

func (p *OutboundProcessor) Init(endpoint ServiceEndpoint, config *OutboundConfig, enableBounce bool) {

}

func (p *OutboundProcessor) SetConnEventHandler(hdlr IConnEventHandler) {
	p.connEvHdlr = hdlr
}
