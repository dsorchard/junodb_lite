package io

import (
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
)

func NewListenerWithFd(cfg ListenerConfig, iocfg InboundConfig, f *os.File, reqHandler IRequestHandler) (lsnr IListener, err error) {
	return
}

func NewListener(cfg ListenerConfig, iocfg InboundConfig, reqHandler IRequestHandler) (lsnr IListener, err error) {
	return
}
