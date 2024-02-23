package io

import "sync"

type InboundConnManager struct {
	mtx         sync.Mutex
	activeConns map[*Connector]struct{}
	wg          sync.WaitGroup
}
