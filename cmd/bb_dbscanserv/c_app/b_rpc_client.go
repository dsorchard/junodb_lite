package app

import "net/rpc"

type RpcClient struct {
	zoneidHash uint32
	timeout    int
	iterations int
	serverIP   string // ip:port
	client     *rpc.Client
}
