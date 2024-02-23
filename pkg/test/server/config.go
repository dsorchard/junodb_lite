package server

import (
	"junodb_lite/cmd/group1/a_proxy/config"
	"junodb_lite/pkg/z_conn_mgr"
)

type ClusterConfig struct {
	//ProxyToBeReplicate      z_conn_mgr.ServiceEndpoint
	//Proxy                   ServerDef
	//StorageServer           ServerDef
	//CAL                     cal.Config
	//Sec                     sec.Config
	ProxyAddress z_conn_mgr.ServiceEndpoint
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
