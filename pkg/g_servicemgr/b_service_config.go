package service

import (
	"junodb_lite/pkg/y_conn_mgr"
	"time"
)

type Config struct {
	Listener            []io.ListenerConfig
	ShutdownWaitTime    time.Duration
	ThrottlingDelayTime time.Duration
	IO                  io.InboundConfigMap
}
