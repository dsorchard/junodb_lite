package z_io

import "net"

type Conn interface {
	GetStateString() string
	GetTLSVersion() string
	GetCipherName() string
	DidResume() string
	GetNetConn() net.Conn
	IsTLS() bool
}
