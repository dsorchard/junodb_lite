package cluster

import (
	"github.com/golang/glog"
)

type Cluster struct {
	Config
	Zones            []*Zone
	RedistSingleZone bool // commit redist one zone only
	RedistZoneId     int  // zone selected for commit one zone.
}

func (c *Cluster) IsRedistZone(zoneid int) bool {
	if !c.RedistSingleZone || (c.RedistZoneId == zoneid) {
		return true
	}
	return false
}

func (c *Cluster) PopulateFromConfig() (err error) {
	return nil

}

func (c *Cluster) WriteToCache(cachePath string, cacheName string, version uint32, forRedist bool) (err error) {
	return err
}

func (c *Cluster) Log() {
	glog.Infof("num of shards: %d, num of zones: %d", c.NumShards, c.NumZones)
	glog.Infof("connInfo: %v", c.ConnInfo)
	for i := uint32(0); i < c.NumZones; i++ {
		c.Zones[i].Log()
	}
	glog.Flush()
}
