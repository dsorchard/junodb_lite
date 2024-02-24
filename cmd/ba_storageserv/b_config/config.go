package config

import (
	"fmt"
	"github.com/golang/glog"
	"junodb_lite/cmd/group1/ba_storageserv/c_storage/db"
	redist "junodb_lite/cmd/group1/ba_storageserv/d_redist"
	config "junodb_lite/cmd/group1/bb_dbscanserv/b_config"
	cluster "junodb_lite/pkg/b_cluster"
	shard "junodb_lite/pkg/b_shard"
	etcd "junodb_lite/pkg/c_etcd"
	initmgr "junodb_lite/pkg/e_initmgr"
	service "junodb_lite/pkg/g_service_mgr"
	io "junodb_lite/pkg/y_conn_mgr"
	"strings"
	"time"
)

type Config struct {
	service.Config
	RootDir     string
	StateLogDir string
	PidFileName string
	LogLevel    string
	ClusterName string
	HttpMonAddr string

	StateLogEnabled    bool
	EtcdEnabled        bool
	MicroShardsEnabled bool
	ShardIdValidation  bool
	CloudEnabled       bool
	DbWatchEnabled     bool

	MaxConcurrentRequests uint32
	NumPrefixDbs          uint32
	NumMicroShards        uint32
	NumMicroShardGroups   uint32
	ReqProcCtxPoolSize    uint32
	MaxTimeToLive         uint32

	RecLockExpiration   time.Duration
	ClusterInfo         *cluster.Config
	DB                  *db.Config
	Redist              *redist.Config
	Etcd                etcd.Config
	ShardMapUpdateDelay time.Duration
	DbScan              config.DbScan

	//DbScan              dbscan.DbScan
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

// Calculate & set the derived info for once
func (c *Config) NewShardMap(zoneId int, machineId int) (shardMap shard.Map) {

	// TODO: validate the rackid and machineid
	node := cluster.ClusterInfo[0].Zones[zoneId].Nodes[machineId]
	shards := node.GetShards()

	//	if err != nil {
	//		glog.Fatalf("Error getting buckets from shard map:%s", err)
	//	}

	shardMap = shard.NewMapWithSize(len(shards))
	for _, s := range shards {
		shardMap[shard.ID(s)] = struct{}{}
	}
	glog.Infof("ShardMap: %v", shardMap)
	return
}

var serverConfig = Config{
	Config: service.Config{
		ShutdownWaitTime: 1 * time.Second,
		IO:               io.InboundConfigMap{},
	},

	LogLevel:              "info",
	PidFileName:           "ss.pid",
	ClusterName:           "cluster",
	EtcdEnabled:           false,
	MaxConcurrentRequests: 3000,
	NumPrefixDbs:          1,
	RecLockExpiration:     600 * time.Millisecond,

	ClusterInfo: &cluster.ClusterInfo[0].Config,

	MicroShardsEnabled:  true,
	NumMicroShards:      0,
	NumMicroShardGroups: 0,

	DB:     &db.DBConfig,
	Redist: &redist.RedistConfig,

	Etcd:                *etcd.NewConfig("127.0.0.1:2379"),
	ShardIdValidation:   true,
	ShardMapUpdateDelay: 30 * time.Second, // 30 seconds
	ReqProcCtxPoolSize:  10000,
	MaxTimeToLive:       3600 * 24 * 3,
}

func ServerConfig() *Config {
	return &serverConfig
}

var Initializer initmgr.IInitializer = initmgr.NewInitializer(initialize, finalize)

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

func LoadConfig(ssConfigFile string) (err error) {
	return
}

func finalize() {
	if serverConfig.EtcdEnabled {
		etcd.Close()
	}
}
