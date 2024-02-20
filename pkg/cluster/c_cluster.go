package cluster

import "errors"

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

type IReader interface {
	// for proxy
	Read(c *Cluster) (version uint32, err error)

	// for storage server
	ReadWithRedistInfo(c *Cluster) (version uint32, err error)

	// for cluster manager
	ReadWithRedistNodeShards(c *Cluster) (err error)
}
