package server

import (
	"junodb_lite/cmd/group1/a_proxy/config"
	"junodb_lite/pkg/io"
)

type ClusterConfig struct {
	//ProxyToBeReplicate      io.ServiceEndpoint
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
