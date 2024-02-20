package cluster

import (
	"fmt"
	"github.com/golang/glog"
	"junodb_lite/pkg/io"
)

var (
	ClusterInfo [2]Cluster
)

func Initialize(args ...interface{}) (err error) {
	//	ccfg *cluster.Config, conf *ioconfig.OutboundConfig, chWatch, etcdReader, cacheFile, cfg.ClusterStats) {
	sz := len(args)
	if sz < 6 {
		err = fmt.Errorf("six arguments expected")
		glog.Error(err)
		return
	}
	var ccfg *Cluster
	var ok bool
	var iocfg *io.OutboundConfig
	var statscfg *StatsConfig

	ccfg, ok = args[0].(*Cluster)
	if !ok {
		err = fmt.Errorf("wrong type of the first argument")
		glog.Error(err)
		return
	}
	iocfg, ok = args[1].(*io.OutboundConfig)
	if !ok {
		err = fmt.Errorf("wrong type of the second argument")
		glog.Error(err)
		return
	}

	if args[5] != nil {
		statscfg, ok = args[5].(*StatsConfig)
		if !ok {
			statscfg = nil
			glog.Exitln("wrong type of the sixth argument")
		}
	}

	if err = InitShardMgr(ccfg, iocfg, statscfg); err != nil {
		glog.Error(err)
		return
	}

	if args[2] != nil && args[3] != nil {
		etcdReader, ok = args[3].(IReader)
		if !ok {
			glog.Exitln("wrong type of the third argument")
		}

		if len(args) >= 5 && args[4] != nil {
			cacheFile, ok = args[4].(string)
			if !ok {
				glog.Exitln("wrong type of the fourth argument")
			}
		}
		var chWatch chan int
		chWatch, ok = args[2].(chan int)
		if ok {
			if !isWatching {
				isWatching = true
				go watchAndResetShardMgr(iocfg, chWatch, statscfg)
			}
		}
	}

	return nil
}
