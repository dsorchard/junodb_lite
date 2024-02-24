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

func (cfg *Config) GetIoConfig(lsnr *io.ListenerConfig) io.InboundConfig {
	if lsnr != nil {

		if c, ok := cfg.IO[lsnr.Name]; ok {
			return c
		} else {
			//if c, ok = cfg.IO[DefaultListenerName]; ok {
			//	return c
			//}
		}
	}
	return io.DefaultInboundConfig
}
