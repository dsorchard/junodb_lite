package io

import (
	"bufio"
	"context"
	util "junodb_lite/pkg/y_util"
	"net"
	"sync"
)

type Connector struct {
	conn net.Conn

	reader        *bufio.Reader
	ctx           context.Context
	cancelCtx     context.CancelFunc
	chResponse    chan IResponseContext
	chStop        chan struct{}
	stopOnce      sync.Once
	closeOnce     sync.Once
	config        InboundConfig
	pendingReq    int32               // local automic counter, for graceful shutdown
	reqCounter    *util.AtomicCounter // global counter
	reqHandler    IRequestHandler
	connMgr       *InboundConnManager
	reqCtxCreator InboundRequestContextCreator
	lsnrType      ListenerType
}

func (c *Connector) Start() {
	go c.doRead()
	go c.doWrite()
}

func (c *Connector) doRead() {
	var magic []byte
	var err error

	r, err := c.reqCtxCreator(magic, c)
	if err != nil {
		return
	}
	r.SetTimeout(c.ctx, c.config.RequestTimeout)
	go c.reqHandler.Process(r)
}

func (c *Connector) doWrite() {

}
