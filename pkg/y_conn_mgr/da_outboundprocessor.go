package io

import (
	"errors"
	util "junodb_lite/pkg/y_util"
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
		shutdown   bool
		numConns   int32 // size of connector pool
		connectors []*OutboundConnector
		doneCh     chan struct{}
		connCh     chan *OutboundConnector
		reqCh      chan IRequestContext
		config     *OutboundConfig
		connInfo   ServiceEndpoint
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

func (p *OutboundProcessor) Shutdown() {

	p.shutdown = true
	for i := 0; i < int(p.numConns); i++ {
		if p.connectors[i] != nil {
			p.connectors[i].Shutdown()
		}
	}
	close(p.doneCh)
}

func (p *OutboundProcessor) WaitShutdown() {
	p.wg.Wait()
	close(p.connCh)
	close(p.reqCh)
}

func (p *OutboundProcessor) SendRequest(req IRequestContext) (err error) {
	return p.sendRequest(req)
}
func (p *OutboundProcessor) sendRequest(req IRequestContext) (err error) {
	// send request
	select {
	case p.reqCh <- req:
	default:
		return errors.New("busy")
	}
	return nil
}
func (p *OutboundProcessor) connect(connCh chan *OutboundConnector, id int, connector *OutboundConnector) {

	interval := p.config.ReconnectIntervalBase
	timer := util.NewTimerWrapper(time.Duration(interval) * time.Millisecond)

	for {
		if p.shutdown {
			return
		}

		select {
		case <-p.doneCh:
			return

		case now := <-timer.GetTimeoutCh():
			conn, err := ConnectTo(&p.connInfo, p.config.ConnectTimeout)
			timeTaken := time.Since(now)
			if err == nil {
				connector = NewOutboundConnector(id, conn.GetNetConn(), p.reqCh, nil, p.config)
				if p.connEvHdlr != nil {
					p.connEvHdlr.OnConnectSuccess(conn, connector, timeTaken)
				}

			}
		}
	}
}

func NewOutboundProcessor(connInfo string, config *OutboundConfig, enableBounce bool) *OutboundProcessor {
	return NewOutbProcessor(ServiceEndpoint{Addr: connInfo}, config, enableBounce)
}

func NewOutbProcessor(endpoint ServiceEndpoint, config *OutboundConfig, enableBounce bool) (p *OutboundProcessor) {
	p = &OutboundProcessor{}
	p.Init(endpoint, config, enableBounce)

	p.Start()
	return p
}
