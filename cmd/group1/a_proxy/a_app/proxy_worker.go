package app

import (
	"github.com/golang/glog"
	config "junodb_lite/cmd/group1/a_proxy/b_config"
	cluster "junodb_lite/pkg/b_cluster"
	etcd "junodb_lite/pkg/c_etcd"
	initmgr "junodb_lite/pkg/e_initmgr"
	util "junodb_lite/pkg/y_util"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type (
	Worker struct {
		CmdProxyCommon
		optWorkerId        uint
		optListenAddresses util.StringListFlags
		optIsChild         bool
		optHttpMonAddr     string
	}
	acceptLimiterT struct {
		acceptDelayTime time.Duration
	}
)

func (c *Worker) GetName() string {
	//TODO implement me
	panic("implement me")
}

func (c *Worker) GetDesc() string {
	//TODO implement me
	panic("implement me")
}

func (c *Worker) GetSynopsis() string {
	//TODO implement me
	panic("implement me")
}

func (c *Worker) GetDetails() string {
	//TODO implement me
	panic("implement me")
}

func (c *Worker) GetOptionDesc() string {
	//TODO implement me
	panic("implement me")
}

func (c *Worker) GetExample() string {
	//TODO implement me
	panic("implement me")
}

func (c *Worker) AddExample(cmdExample string, desc string) {
	//TODO implement me
	panic("implement me")
}

func (c *Worker) AddDetails(txt string) {
	//TODO implement me
	panic("implement me")
}

func (c *Worker) Exec() {
	numInheritedFDs := util.GetNumOpenFDs()

	cfg := &config.Conf

	initmgr.Register(config.Initializer, c.optConfigFile)
	initmgr.Init()

	if len(c.optListenAddresses) != 0 {
		cfg.SetListeners(c.optListenAddresses)
	}

	var chWatch chan int
	var etcdReader cluster.IReader
	if cfg.EtcdEnabled {
		chWatch = etcd.WatchForProxy()
		etcdReader = etcd.GetClsReadWriter()
	}
	cacheFile := filepath.Join(cfg.Etcd.CacheDir, cfg.Etcd.CacheName)
	initmgr.RegisterWithFuncs(cluster.Initialize, cluster.Finalize, &cluster.ClusterInfo[0],
		&cfg.Outbound, chWatch, etcdReader, cacheFile)

	initmgr.RegisterWithFuncs(replication.Initialize, replication.Finalize, &cfg.Replication)
	if cfg.EtcdEnabled {
		initmgr.RegisterWithFuncs(watcher.Initialize, watcher.Finalize, cfg.ClusterName, etcd.GetEtcdCli(), &cfg.Etcd)
	}

	initmgr.Init()

	var service *service.Service
	httpEnabled := len(cfg.HttpMonAddr) != 0

	if c.optIsChild {
		numListeners := len(cfg.Listener)

		if numInheritedFDs >= numListeners+3 {
			var fds []*os.File

			for i := 0; i < numListeners; i++ {
				if f := os.NewFile(uintptr(3+i), ""); f != nil && util.IsSocket(f) {
					fds = append(fds, f)
				} else {
					glog.Exitf("fd not validate")
				}
			}
			if httpEnabled {
				if numInheritedFDs > numListeners+3 {
					f := os.NewFile(uintptr(numListeners+3), "")
					if f != nil && util.IsSocket(f) {
						if httpListener, err := net.FileListener(f); err == nil {
							shmstats.SetHttpPort(httpListener.Addr().String())
							go func() {
								http.Serve(httpListener, &stats.HttpServerMux)
							}()
						}
					}
				}
			}

			service = handler.NewProxyServiceWithListenFd(&config.Conf, &acceptLimiterT{acceptDelayTime: cfg.Config.ThrottlingDelayTime.Duration}, fds...)
		}

	} else {
		if httpEnabled {
			go func() {
				glog.Infof("to serve HTTP on %s", cfg.HttpMonAddr)
				if err := http.ListenAndServe(cfg.HttpMonAddr, &stats.HttpServerMux); err != nil {
					glog.Warningf("fail to serve HTTP on %s, err: %s", cfg.HttpMonAddr, err)
				}
			}()

		}
		service = handler.NewProxyService(cfg)
	}
	service.Run()
}

func (c *Worker) PrintUsage() {
	//TODO implement me
	panic("implement me")
}

func (c *Worker) Parse(args []string) error {
	//TODO implement me
	panic("implement me")
}

func (c *Worker) Init(name string, desc string) {
	//c.CmdProxyCommon.Init(name, desc)
	//c.UintOption(&c.optWorkerId, "id|worker-id", 0, "specify the ID of the worker")
	//c.ValueOption(&c.optListenAddresses, "listen", "specify listening address. Override Listener in config file")
	//c.BoolOption(&c.optIsChild, "child", false, "specify if the worker was started by a parent process")
	//c.StringOption(&c.optHttpMonAddr, "mon-addr|monitoring-address", "", "specify the http monitoring address. \n\toverride HttpMonAddr in config file")
}
