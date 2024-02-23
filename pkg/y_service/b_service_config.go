package service

import (
	"junodb_lite/pkg/io"
	"time"
)

type Config struct {
	Listener            []io.ListenerConfig
	ShutdownWaitTime    time.Duration
	ThrottlingDelayTime time.Duration
	IO                  io.InboundConfigMap
}
