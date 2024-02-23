package z_conn_mgr

import "net"

type Conn interface {
	GetStateString() string
	GetTLSVersion() string
	GetCipherName() string
	DidResume() string
	GetNetConn() net.Conn
	IsTLS() bool
}
