package io

import (
	"context"
	"net"
	"os"
	"time"
)

type (
	ListenerType byte

	IListener interface {
		GetName() string
		GetType() ListenerType
		AcceptAndServe() error
		Close() error
		Shutdown()
		WaitForShutdownToComplete(time.Duration)
		GetConnString() string
		GetNumActiveConnections() uint32
		Refresh()
	}
	Listener struct {
		config      ListenerConfig
		ioConfig    InboundConfig
		netListener net.Listener
		reqHandler  IRequestHandler
		connMgr     *InboundConnManager
		lsnrType    ListenerType
	}
)

func NewListenerWithFd(cfg ListenerConfig, iocfg InboundConfig, f *os.File, reqHandler IRequestHandler) (lsnr IListener, err error) {
	return
}

func NewListener(cfg ListenerConfig, iocfg InboundConfig, reqHandler IRequestHandler) (lsnr IListener, err error) {
	return
}

var _ IListener = new(Listener)

func (l *Listener) Close() error {
	return l.netListener.Close()
}
func (l *Listener) GetName() string {
	//TODO implement me
	panic("implement me")
}

func (l *Listener) GetType() ListenerType {
	//TODO implement me
	panic("implement me")
}

func (l *Listener) AcceptAndServe() error {
	conn, err := l.netListener.Accept()

	if err == nil {
		//if cal.IsEnabled() {
		//	raddr := conn.RemoteAddr().String()
		//	if rhost, _, e := net.SplitHostPort(raddr); e == nil {
		//		cal.Event(cal.TxnTypeAccept, rhost, cal.StatusSuccess, []byte("raddr="+raddr+"&laddr="+conn.LocalAddr().String()))
		//	}
		//}
		//otel.RecordCount(otel.Accept, []otel.Tags{{otel.Status, otel.Success}})
		l.startNewConnector(conn)
	} else {
		//otel.RecordCount(otel.Accept, []otel.Tags{{otel.Status, otel.Error}})
	}
	//log the error in caller if needed
	return err
}

func (l *Listener) startNewConnector(conn net.Conn) {
	ctx, cancel := context.WithCancel(context.Background())
	bufSize := l.ioConfig.IOBufSize
	if bufSize == 0 {
		bufSize = 64000
	}
	connector := &Connector{
		conn: conn,
		//reader:     util.NewBufioReader(conn, bufSize),
		ctx:        ctx,
		cancelCtx:  cancel,
		chResponse: make(chan IResponseContext, l.ioConfig.RespChanSize),
		chStop:     make(chan struct{}),
		connMgr:    l.connMgr,
		reqHandler: l.reqHandler,
		config:     l.ioConfig,
		pendingReq: 0,
		//		reqCounter:    server.GetReqCounter(),
		reqCtxCreator: l.reqHandler.GetReqCtxCreator(),
		lsnrType:      l.GetType(),
	}
	//if connector.reqCtxCreator == nil {
	//	connector.reqCtxCreator = DefaultInboundRequestContexCreator
	//}
	connector.Start()
}

func (l *Listener) Shutdown() {
	//TODO implement me
	panic("implement me")
}

func (l *Listener) WaitForShutdownToComplete(duration time.Duration) {
	//TODO implement me
	panic("implement me")
}

func (l *Listener) GetConnString() string {
	//TODO implement me
	panic("implement me")
}

func (l *Listener) GetNumActiveConnections() uint32 {
	//TODO implement me
	panic("implement me")
}

func (l *Listener) Refresh() {
	//TODO implement me
	panic("implement me")
}
