package server

import (
	"junodb_lite/cmd/a_proxy/b_config"
	"junodb_lite/pkg/y_conn_mgr"
)

type ClusterConfig struct {
	//ProxyToBeReplicate      y_conn_mgr.ServiceEndpoint
	//Proxy                   ServerDef
	//StorageServer           ServerDef
	//CAL                     cal.Config
	//Sec                     sec.Config
	ProxyAddress io.ServiceEndpoint
	ProxyConfig  *config.Config

	SSdir                   string
	WalDir                  string
	LogLevel                string
	Proxydir                string
	Githubdir               string
	MarkDown                string
	RedistType              string
	SecondHostSSdir         string
	AddRemoveSecondHost     string
	EtcdServerRestartGitDir string
}
