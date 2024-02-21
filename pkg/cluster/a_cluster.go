package cluster

import (
	"errors"
	"github.com/golang/glog"
)

type IReader interface {
	// for proxy
	Read(c *Cluster) (version uint32, err error)

	// for storage server
	ReadWithRedistInfo(c *Cluster) (version uint32, err error)

	// for cluster manager
	ReadWithRedistNodeShards(c *Cluster) (err error)
}

type Cluster struct {
	Config
	Zones            []*Zone
	RedistSingleZone bool // commit redist one zone only
	RedistZoneId     int  // zone selected for commit one zone.
}

func (c *Cluster) Read(r IReader) (version uint32, err error) {
	if r == nil {
		return 0, errors.New("nil cluster reader")
	}

	return r.Read(c)
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
