package service

import (
	"junodb_lite/pkg/z_conn_mgr"
	"time"
)

type Config struct {
	Listener            []z_conn_mgr.ListenerConfig
	ShutdownWaitTime    time.Duration
	ThrottlingDelayTime time.Duration
	IO                  z_conn_mgr.InboundConfigMap
}
