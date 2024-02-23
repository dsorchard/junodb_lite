package service

import (
	"junodb_lite/pkg/z_io"
	"time"
)

type Config struct {
	Listener            []z_io.ListenerConfig
	ShutdownWaitTime    time.Duration
	ThrottlingDelayTime time.Duration
	IO                  z_io.InboundConfigMap
}
