package redist

import "junodb_lite/pkg/y_conn_mgr"

type Config struct {
	Outbound io.OutboundConfig
}

var RedistConfig = DefRedistConfig

var DefRedistConfig = Config{}
