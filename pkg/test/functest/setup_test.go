package functest

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/golang/glog"
	"junodb_lite/cmd/group1/a_proxy/config"
	"junodb_lite/pkg/b_cluster"
	"junodb_lite/pkg/c_etcd"
	"junodb_lite/pkg/test/server"
	"junodb_lite/pkg/z_io"
	"net/http"
	"os"
	"testing"
)

var testConfig = server.ClusterConfig{
	ProxyAddress: z_io.ServiceEndpoint{
		Addr:       "127.0.0.1:26969",
		SSLEnabled: false,
	},
	ProxyConfig: &config.Conf,
	//Sec:         sec.DefaultConfig,
}

func TestMain(m *testing.M) {
	var (
		configFile  string
		logLevel    string
		longevity   string
		httpEnabled bool
	)
	flag.StringVar(&logLevel, "log_level", "", "specify log level")
	flag.StringVar(&configFile, "config", "", "specify config file")
	flag.StringVar(&longevity, "long", "false", "specify longevity test or not")
	flag.BoolVar(&httpEnabled, "http", true, "enable http")
	flag.Parse()
	if len(configFile) == 0 {
		os.Exit(-1)
	}
	if httpEnabled {
		go func() {
			http.ListenAndServe(":6060", nil)
		}()
	}

	testConfig.LogLevel = "warning"

	if _, err := toml.DecodeFile(configFile, &testConfig); err != nil {
		fmt.Printf("fail to read %s. error: %s", configFile, err)
		os.Exit(-1)
	}
	if logLevel != "" {
		testConfig.LogLevel = logLevel
	}
	os.Unsetenv("JUNO_PIN")

	var chWatch chan int
	clusterInfo := &cluster.ClusterInfo[0]
	glog.Info("preCluster info: ", clusterInfo)
	if testConfig.ProxyConfig.EtcdEnabled {
		chWatch = etcd.WatchForProxy()
		etcd.Connect(&testConfig.ProxyConfig.Etcd, testConfig.ProxyConfig.ClusterName)
		rw := etcd.GetClsReadWriter()
		if rw == nil {
			glog.Exitf("no etcd setup")
		}
		clusterInfo.Read(rw)
	} else {
		clusterInfo.PopulateFromConfig()
	}
	cluster.Initialize(&cluster.ClusterInfo[0], &testConfig.ProxyConfig.Outbound, chWatch, etcd.GetClsReadWriter(), nil, nil)
	glog.Info("postCluster info: ", clusterInfo)

	rc := m.Run()
	os.Exit(rc)
}
