package service

import (
	"junodb_lite/pkg/y_conn_mgr"
	"time"
)

type Config struct {
	Listener            []y_conn_mgr.ListenerConfig
	ShutdownWaitTime    time.Duration
	ThrottlingDelayTime time.Duration
	IO                  y_conn_mgr.InboundConfigMap
}
