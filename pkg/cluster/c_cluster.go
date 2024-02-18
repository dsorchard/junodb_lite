package cluster

type Cluster struct {
	Config
	Zones            []*Zone
	RedistSingleZone bool // commit redist one zone only
	RedistZoneId     int  // zone selected for commit one zone.
}

type IReader interface {
	// for proxy
	Read(c *Cluster) (version uint32, err error)

	// for storage server
	ReadWithRedistInfo(c *Cluster) (version uint32, err error)

	// for cluster manager
	ReadWithRedistNodeShards(c *Cluster) (err error)
}
