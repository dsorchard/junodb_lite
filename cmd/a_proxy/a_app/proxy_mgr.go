package app

import (
	"fmt"
	"github.com/golang/glog"
	"junodb_lite/cmd/group1/a_proxy/b_config"
	initmgr "junodb_lite/pkg/e_initmgr"
	sec "junodb_lite/pkg/ea_sec"
	service "junodb_lite/pkg/g_service_mgr"
	util "junodb_lite/pkg/y_util"
	"os"
	"strconv"
	"strings"
	"syscall"
)

type (
	Manager struct {
		CmdProxyCommon
		optNumChildren     uint
		optListenAddresses util.StringListFlags
		optHttpMonAddr     string
		cmdArgs            []string
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

	cfg := &config.Conf

	pidFile := cfg.PidFileName

	if data, err := os.ReadFile(pidFile); err == nil {
		if pid, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil {
			if process, err := os.FindProcess(pid); err == nil {
				if err := process.Signal(syscall.Signal(0)); err == nil {
					glog.Exitf("process pid: %d in %s is still running\n", pid, pidFile)
					///TODO check if it is proxy process
				}
			}
		}
	}
	os.WriteFile(pidFile, []byte(fmt.Sprintf("%d\n", os.Getpid())), 0644)
	defer os.Remove(pidFile)
	initmgr.Register(sec.Initializer, true)

	initmgr.Init()

	servermgr := service.NewServerManager(int(cfg.NumChildren), os.Args[0], c.cmdArgs,
		cfg.Config, cfg.HttpMonAddr, cfg.CloudEnabled)
	servermgr.Run()
}

func (c *Manager) PrintUsage() {
	//TODO implement me
	panic("implement me")
}
func (c *Manager) Init(name string, desc string) {
	//c.CmdProxyCommon.Init(name, desc)
	//c.UintOption(&c.optNumChildren, "n|num-children", kDefaultNumChild, "specify the number of worker process(es)")
	//c.FlagSet.Var(&c.optListenAddresses, "listen", "specify listening address")
	//c.StringOption(&c.optHttpMonAddr, "mon-addr|monitoring-address", "", "specify the http monitoring address. \n\toverride HttpMonAddr in config file")
}
