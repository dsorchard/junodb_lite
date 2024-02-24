package app

import (
	"github.com/golang/glog"
	stats "junodb_lite/cmd/a_proxy/e_stats"
	handler "junodb_lite/cmd/a_proxy/fa_handler"
	config "junodb_lite/cmd/ba_storageserv/b_config"
	storage "junodb_lite/cmd/ba_storageserv/c_storage"
	redist "junodb_lite/cmd/ba_storageserv/d_redist"
	watcher "junodb_lite/cmd/ba_storageserv/e_watcher"
	compact "junodb_lite/cmd/ba_storageserv/f_compact"
	patch "junodb_lite/cmd/bb_dbscanserv/a_patch"
	initmgr "junodb_lite/pkg/e_initmgr"
	service "junodb_lite/pkg/g_service_mgr"
	util "junodb_lite/pkg/y_util"
	"net"
	"net/http"
	"os"
	"strconv"
)

type Worker struct {
	CmdStorageCommon
	optWorkerId        uint
	optListenAddresses util.StringListFlags
	optIsChild         bool
	optHttpMonAddr     string
	optZoneId          uint
	optMachineIndex    uint
	optLRUCacheSize    uint
}

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

	initmgr.Register(config.Initializer, c.optConfigFile)
	initmgr.Init() //initalize config first as others depend on it

	cfg := config.ServerConfig()
	if len(c.optListenAddresses) != 0 {
		cfg.SetListeners(c.optListenAddresses)
	}
	if len(c.optHttpMonAddr) != 0 {
		cfg.HttpMonAddr = c.optHttpMonAddr
	}

	if _, err := strconv.Atoi(cfg.HttpMonAddr); err == nil {
		cfg.HttpMonAddr = ":" + cfg.HttpMonAddr
	}

	initmgr.RegisterWithFuncs(storage.Initialize, storage.Finalize, int(c.optZoneId), int(c.optMachineIndex), int(c.optLRUCacheSize))
	initmgr.Init()

	patch.Init(&cfg.DbScan) // for namespace migration
	if cfg.EtcdEnabled {
		watcher.Init(cfg.ClusterName,
			uint16(c.optZoneId),
			uint16(c.optMachineIndex),
			&(cfg.Etcd),
			cfg.ShardMapUpdateDelay,
			1)

		redist.Init(cfg.ClusterName,
			uint16(c.optZoneId),
			uint16(c.optMachineIndex),
			&(cfg.Etcd))
	}

	reqHandler := handler.NewRequestHandler()

	service, suspend := service.NewService(cfg.Config, reqHandler)

	if len(cfg.HttpMonAddr) != 0 {
		if c.optIsChild {
			if numInheritedFDs > 3 {
				if f := os.NewFile(3, ""); f != nil {
					if httpListener, err := net.FileListener(f); err == nil {
						go func() {
							http.Serve(httpListener, &stats.HttpServerMux)
						}()
					}
				}
			} else {
				glog.Warningf("no inherited fds")
			}
		} else {
			go func() {
				glog.Infof("to serve HTTP on %s", cfg.HttpMonAddr)
				if err := http.ListenAndServe(cfg.HttpMonAddr, &stats.HttpServerMux); err != nil {
					glog.Warningf("fail to serve HTTP on %s, err: %s", cfg.HttpMonAddr, err)
				}
			}()
		}
	}

	service.Zoneid = int(c.optZoneId)
	if cfg.DbWatchEnabled {
		go compact.Watch(int(c.optZoneId), int(c.optMachineIndex), suspend)
	}
	service.Run()
}

func (c *Worker) PrintUsage() {
	//TODO implement me
	panic("implement me")
}

func (c *Worker) Init(name string, desc string) {
	//c.CmdStorageCommon.Init(name, desc)
	//c.UintOption(&c.optWorkerId, "id|worker-id", 0, "specify the ID of the worker")
	//c.ValueOption(&c.optListenAddresses, "listen", "specify listening address. Override Listener in config file")
	//c.BoolOption(&c.optIsChild, "child", false, "specify if the worker was started by a parent process")
	//c.StringOption(&c.optHttpMonAddr, "mon-addr|monitoring-address", "", "specify the http monitoring address. \n\toverride HttpMonAddr in config file")
	//c.UintOption(&c.optZoneId, "zone-id", 0, "specify zone id")
	//c.UintOption(&c.optMachineIndex, "machine-index", 0, "specify machine index")
	//c.UintOption(&c.optLRUCacheSize, "lru-cache-mb", 0, "specify lru cache size")
}
