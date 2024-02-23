package y_conn_mgr

import "sync"

type InboundConnManager struct {
	mtx         sync.Mutex
	activeConns map[*Connector]struct{}
	wg          sync.WaitGroup
}
