package io

import (
	"net"
	"time"
)

type Conn interface {
	GetStateString() string
	GetTLSVersion() string
	GetCipherName() string
	DidResume() string
	GetNetConn() net.Conn
	IsTLS() bool
}

func ConnectTo(endpoint *ServiceEndpoint, connectTimeout time.Duration) (conn Conn, err error) {
	return nil, nil
}
