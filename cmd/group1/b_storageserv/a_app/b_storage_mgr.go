package app

import (
	"fmt"
	"github.com/golang/glog"
	config "junodb_lite/cmd/group1/b_storageserv/b_config"
	initmgr "junodb_lite/pkg/e_initmgr"
	"net"
	"os"
)

type (
	Manager struct {
		CmdStorageCommon
		optNumChildren uint
		optIpAddress   string
		cmdArgs        []string
	}
)

func (c *Manager) GetName() string {
	//TODO implement me
	panic("implement me")
}

func (c *Manager) GetDesc() string {
	//TODO implement me
	panic("implement me")
}

func (c *Manager) GetSynopsis() string {
	//TODO implement me
	panic("implement me")
}

func (c *Manager) GetDetails() string {
	//TODO implement me
	panic("implement me")
}

func (c *Manager) GetOptionDesc() string {
	//TODO implement me
	panic("implement me")
}

func (c *Manager) GetExample() string {
	//TODO implement me
	panic("implement me")
}

func (c *Manager) AddExample(cmdExample string, desc string) {
	//TODO implement me
	panic("implement me")
}

func (c *Manager) AddDetails(txt string) {
	//TODO implement me
	panic("implement me")
}

func (c *Manager) Exec() {
	initmgr.Register(config.Initializer, c.optConfigFile)
	initmgr.Init() //initalize config first as others depend on it

	cfg := config.ServerConfig()

	var connInfo []ConnectInfo
	if c.optIpAddress == "" {
		for row := 0; row < len(cfg.ClusterInfo.ConnInfo); row++ {
			for col := 0; col < len(cfg.ClusterInfo.ConnInfo[row]); col++ {
				ipport := cfg.ClusterInfo.ConnInfo[row][col]
				if ip, _, err := net.SplitHostPort(ipport); err == nil {
					{
						// for K8s storage pod initialization check in GKE
						k8sPodName, ok := os.LookupEnv("POD_NAME")
						if ok {
							k8sPodFqdn := k8sPodName
							k8sPodDomain, ok := os.LookupEnv("POD_DOMAIN")
							if ok {
								k8sPodFqdn = k8sPodName + "." + k8sPodDomain
							}
							if ip == k8sPodFqdn {
								connInfo = append(connInfo, ConnectInfo{Listener: ipport, ZoneId: row, MachineIndex: col})
							} else {
								glog.Errorf("[K8s] ConnInfo Ip of pod (%c) doesn't match pod fqdn (%c)", ip, k8sPodFqdn)
							}
						}
					}
				} else {
					glog.Errorf("wrong connect info string %c [%d][%d]", ipport, row, col)
				}

			}
		}
	} else {
		for row := 0; row < len(cfg.ClusterInfo.ConnInfo); row++ {
			for col := 0; col < len(cfg.ClusterInfo.ConnInfo[row]); col++ {
				ipport := cfg.ClusterInfo.ConnInfo[row][col]
				if ip, _, err := net.SplitHostPort(ipport); err == nil {
					if ip == c.optIpAddress {
						connInfo = append(connInfo, ConnectInfo{Listener: ipport, ZoneId: row, MachineIndex: col})
					}
				} else {
					glog.Errorf("wrong connect info string %c [%d][%d]", ipport, row, col)
				}
			}
		}
	}
	if len(connInfo) == 0 {
		glog.Errorf("No cluster info found for this host, exit")
		return
	}

	initmgr.Init()

	// Parent process
	var cmdArgs []string

	if c.optConfigFile != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("-config=%c", c.optConfigFile))
	}
	servermgr := NewServerManager(len(connInfo), cfg.PidFileName, os.Args[0], cmdArgs, connInfo,
		cfg.HttpMonAddr, int(8080), cfg.CloudEnabled)
	servermgr.Run()
}

func (c *Manager) PrintUsage() {
	//TODO implement me
	panic("implement me")
}

func (c *Manager) Init(name string, desc string) {
}
