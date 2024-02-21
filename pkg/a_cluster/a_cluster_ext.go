package cluster

type IReader interface {
	Read(c *Cluster) (version uint32, err error)               // for proxy
	ReadWithRedistInfo(c *Cluster) (version uint32, err error) // for storage server
	ReadWithRedistNodeShards(c *Cluster) (err error)           // for cluster manager
}
