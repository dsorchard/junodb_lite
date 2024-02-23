package server

import (
	"junodb_lite/cmd/group1/a_proxy/config"
	"junodb_lite/pkg/z_io"
)

type ClusterConfig struct {
	//ProxyToBeReplicate      z_io.ServiceEndpoint
	//Proxy                   ServerDef
	//StorageServer           ServerDef
	//CAL                     cal.Config
	//Sec                     sec.Config
	ProxyAddress z_io.ServiceEndpoint
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
