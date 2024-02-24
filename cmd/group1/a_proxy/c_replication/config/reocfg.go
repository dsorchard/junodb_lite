package config

import io "junodb_lite/pkg/y_conn_mgr"

type (
	ReplicationTarget struct {
		Name string
		io.ServiceEndpoint
		UseMayflyProtocol bool
		Namespaces        []string
		BypassLTMEnabled  bool
	}

	Config struct {
		Targets []ReplicationTarget
		IO      io.OutboundConfigMap
	}
)
