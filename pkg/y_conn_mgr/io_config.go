package io

import "time"

type (
	InboundConfig struct {
		HandshakeTimeout     time.Duration //only for TLS connection
		IdleTimeout          time.Duration
		ReadTimeout          time.Duration
		WriteTimeout         time.Duration
		RequestTimeout       time.Duration
		MaxBufferedWriteSize int
		IOBufSize            int
		RespChanSize         int
	}
	OutboundConfig struct {
		ReconnectIntervalBase int
		ConnectTimeout        time.Duration
	}
	InboundConfigMap  map[string]InboundConfig
	OutboundConfigMap map[string]OutboundConfig
)

var (
	DefaultInboundConfig = InboundConfig{
		HandshakeTimeout:     500 * time.Millisecond,
		IdleTimeout:          120 * time.Second,
		ReadTimeout:          500 * time.Millisecond,
		WriteTimeout:         500 * time.Millisecond,
		RequestTimeout:       600 * time.Millisecond,
		MaxBufferedWriteSize: 64 * 1024,
		IOBufSize:            64 * 1024, // default 64k buf size
		RespChanSize:         10000,
	}
)
