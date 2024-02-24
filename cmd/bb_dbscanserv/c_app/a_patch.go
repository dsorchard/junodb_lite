package app

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/tecbot/gorocksdb"
	config "junodb_lite/cmd/bb_dbscanserv/b_config"
	"strings"
	"sync"
	"time"
)

var (
	writeOptions *gorocksdb.WriteOptions = gorocksdb.NewDefaultWriteOptions()
	dbclient     *gorocksdb.DB
	client       *RpcClient

	namespaces string
	nsList     [][]byte

	relayEnabled bool
	patchDbPath  string
	patchTTL     = 86400 * 8 // Default to 8 days
	debug        bool

	rwMutex    sync.RWMutex
	E9         = uint64(time.Second)
	TEST_PATCH = "__test_patch"
)

// Called by storageserv
func InitPatch(cfg *config.DbScan) {
	namespaces = cfg.ReplicationNamespaces
	names := strings.Split(namespaces, "|")
	nsList = make([][]byte, len(names))
	for i := range names {
		nsList[i] = append(nsList[i], names[i]...)
	}

	if len(cfg.ReplicationAddr) == 0 {
		return
	}
	if len(names) > 0 {
		relayEnabled = true
		if cfg.PatchTTL > 0 {
			patchTTL = cfg.PatchTTL
		}
	}

	debug = cfg.Debug
	glog.Infof("Init %v", *cfg)
	client = NewRpcClient(cfg.ListenPort)
}

func NewRpcClient(port int) *RpcClient {
	client := &RpcClient{
		serverIP: fmt.Sprintf(":%d", port),
	}
	client.timeout = 20
	client.iterations = 3
	return client
}
