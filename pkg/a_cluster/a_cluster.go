package cluster

import (
	"errors"
	"github.com/golang/glog"
)

type Cluster struct {
	Config
	Zones            []*Zone
	RedistSingleZone bool // commit redist one zone only
	RedistZoneId     int  // zone selected for commit one zone.

	shardMap *cluster.ShardMap
}

func (c *Cluster) GetShards(zoneid uint32, nodeid uint32) (shards []uint32, err error) {
	if zoneid >= c.NumZones {
		return nil, errors.New("invalid zone id")
	}

	if nodeid >= c.Zones[zoneid].NumNodes {
		return nil, errors.New("invalid node id")
	}

	return c.Zones[zoneid].Nodes[nodeid].GetShards(), nil
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

func (n *Node) allocate(primaryCount int, secondaryCount int) {
	n.PrimaryShards = make([]uint32, 0, primaryCount)
	n.SecondaryShards = make([]uint32, 0, secondaryCount)
}

func (n *Node) fillPrimary(start int, count int, shards []uint32) {
	n.PrimaryShards = append(n.PrimaryShards, shards[start:start+count]...)
}

func (n *Node) fillSecondary(start int, count int, shards []uint32) {
	n.SecondaryShards = append(n.SecondaryShards, shards[start:start+count]...)
}
