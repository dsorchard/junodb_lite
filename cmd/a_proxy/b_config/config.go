package config

import (
	"fmt"
	repconfig "junodb_lite/cmd/a_proxy/c_replication/config"
	"junodb_lite/pkg/c_etcd"
	initmgr "junodb_lite/pkg/e_initmgr"
	service "junodb_lite/pkg/g_service_mgr"
	io "junodb_lite/pkg/y_conn_mgr"
	"strings"
)

type Config struct {
	EtcdEnabled  bool
	Etcd         etcd.Config
	ClusterName  string
	Outbound     interface{}
	PidFileName  string
	NumChildren  int
	Config       service.Config
	HttpMonAddr  string
	CloudEnabled bool
	Listener     []io.ListenerConfig
	Replication  repconfig.Config
}

func (cfg *Config) SetListeners(values []string) {
	cfg.Listener = make([]io.ListenerConfig, len(values))
	for i, str := range values {
		str = strings.ToLower(str)
		lncfg := &cfg.Listener[i]
		if strings.HasPrefix(str, "ssl:") {
			str = strings.TrimPrefix(str, "ssl:")
			lncfg.SSLEnabled = true
		}
		if !strings.Contains(str, ":") {
			lncfg.Addr = ":" + str
		} else {
			lncfg.Addr = str
		}
	}
}

var (
	Initializer initmgr.IInitializer = initmgr.NewInitializer(initialize, finalize)
	Conf                             = Config{}
)

func initialize(args ...interface{}) (err error) {
	sz := len(args)
	if sz < 1 {
		err = fmt.Errorf("a string config file name argument expected")
		return
	}
	filename, ok := args[0].(string)

	if ok == false {
		err = fmt.Errorf("wrong argument type. a string config file name expected")
		return
	}
	err = LoadConfig(filename)
	return
}

func LoadConfig(file string) (err error) {
	return nil
}

func finalize() {
	if Conf.EtcdEnabled {
		etcd.Close()
	}
}
