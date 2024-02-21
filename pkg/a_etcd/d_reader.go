package etcd

import (
	"errors"
	"github.com/golang/glog"
	"junodb_lite/pkg/b_cluster"
	"strconv"
	"strings"
)

// Implements b_cluster.IReader
type Reader struct {
	etcdcli *Client
}

var _ cluster.IReader = new(Reader)

func (cr *Reader) Read(c *cluster.Cluster) (version uint32, err error) {

	var algVer uint32
	if version, algVer, err = cr.etcdcli.GetVersion(); err != nil {
		return
	}

	c.AlgVersion = algVer
	if c.NumShards, err = cr.etcdcli.GetUint32(TagNumShards); err != nil {
		return
	}

	if c.NumZones, err = cr.etcdcli.GetUint32(TagNumZones); err != nil {
		return
	}
	c.ConnInfo = make([][]string, c.NumZones)
	c.Zones = make([]*cluster.Zone, c.NumZones)

	if err = cr.readNodesIpport(c, TagNodeIpport, 2); err != nil {
		return
	}

	if err = cr.readNodesShards(c, TagNodeShards, 2); err != nil {
		return
	}

	for zoneid := 0; zoneid < int(c.NumZones); zoneid++ {
		if c.ConnInfo[zoneid] == nil {
			glog.Errorf("[ERROR]: ip:port is missing for zone %d in etcd.", zoneid)
			return 0, errors.New("Missing ip:port in etcd.")
		}

		if c.Zones[zoneid] == nil {
			glog.Errorf("[ERROR]: shardmap is missing for zone %d in etcd.", zoneid)
			return 0, errors.New("Missing shardmap info in etcd.")
		}
	}

	err1 := c.WriteToCache(cr.etcdcli.config.CacheDir, cr.etcdcli.config.CacheName,
		version, false)
	if err1 != nil {
		glog.Errorf("failed to write to etcd cache: %s", err1.Error())
	}

	return
}

func (cr *Reader) ReadWithRedistInfo(c *cluster.Cluster) (version uint32, err error) {

	// read nodes from cluster info
	if version, err = cr.Read(c); err != nil {
		return
	}

	//	glog.Infof("initial cluster")
	//	c.Log()

	// read nodes from redist
	nc := &cluster.Cluster{
		Config: cluster.Config{
			NumZones:  c.NumZones,
			NumShards: c.NumShards,
		},
	}

	nc.ConnInfo = make([][]string, c.NumZones)
	nc.Zones = make([]*cluster.Zone, c.NumZones)
	if err = cr.readNodesIpport(nc, TagRedistNodeIpport, 3); err != nil {
		return
	}

	if err = cr.readNodesShards(nc, TagRedistNodeShards, 3); err != nil {
		return
	}

	if nc.Zones == nil {
		return
	}

	// merge nodes from redist to cluster info for storage server
	for zoneid := 0; zoneid < int(nc.NumZones); zoneid++ {
		originlen := len(c.ConnInfo[zoneid])
		if len(nc.ConnInfo[zoneid]) > originlen {
			c.ConnInfo[zoneid] = append(c.ConnInfo[zoneid], nc.ConnInfo[zoneid][originlen:]...)
		}
	}

	forRedist := false
	for zoneid := 0; zoneid < int(nc.NumZones); zoneid++ {
		if nc.Zones[zoneid] == nil {
			continue
		}

		originlen := len(c.Zones[zoneid].Nodes)
		if nc.Zones[zoneid] != nil && len(nc.Zones[zoneid].Nodes) > originlen {
			c.Zones[zoneid].Nodes = append(c.Zones[zoneid].Nodes, nc.Zones[zoneid].Nodes[originlen:]...)
			c.Zones[zoneid].NumNodes = uint32(len(c.Zones[zoneid].Nodes))
			forRedist = true
		}
	}
	//	glog.Info("redist cluster")
	//	nc.Log()
	glog.Info("cluster info adjusted with new nodes")
	c.Log()

	err1 := c.WriteToCache(cr.etcdcli.config.CacheDir, cr.etcdcli.config.CacheName,
		version, forRedist)
	if err1 != nil {
		glog.Errorf("failed to write to etcd cache: %s", err1.Error())
	}
	return
}

func (cr *Reader) ReadWithRedistNodeShards(c *cluster.Cluster) (err error) {
	//TODO implement me
	panic("implement me")
}

// tag: either TagNodeIpport or TagRedistNodeIpport
// offset: is the index of the token for zoneid after split with delimiter "_"
// redist_node_ipport_0_1 => 3
// node_ipport_0_0 => 2
func (cr *Reader) readNodesIpport(c *cluster.Cluster, tag string, offset int) (err error) {
	resp, err := cr.etcdcli.getWithPrefix(tag)
	if err != nil {
		return err
	}

	// if resp.Count == 0 {
	//	return errors.New("0 node")
	// }

	for _, ev := range resp.Kvs {
		tokens := strings.Split(string(ev.Key), TagCompDelimiter)
		if len(tokens) < offset+2 {
			// log error?
			continue
		}
		zoneid, _ := strconv.Atoi(tokens[offset])
		nodeid, _ := strconv.Atoi(tokens[offset+1])

		if zoneid >= int(c.NumZones) {
			// log error?
			continue
		}

		// the prefix fetch is sorted by key in reverse order
		if c.ConnInfo[zoneid] == nil {
			c.ConnInfo[zoneid] = make([]string, nodeid+1)
		}

		c.ConnInfo[zoneid][nodeid] = string(ev.Value)
	}
	return nil
}

// tag: either TagNodeShards or TagRedistNodeShards
func (cr *Reader) readNodesShards(c *cluster.Cluster, tag string, offset int) (err error) {
	resp, err := cr.etcdcli.getWithPrefix(tag)
	if err != nil {
		return err
	}

	// if resp.Count == 0 {
	//	return errors.New("0 node")
	// }

	for _, ev := range resp.Kvs {
		tokens := strings.Split(string(ev.Key), TagCompDelimiter)
		if len(tokens) < offset+2 {
			// log error?
			continue
		}
		zoneid, _ := strconv.Atoi(tokens[offset])
		nodeid, _ := strconv.Atoi(tokens[offset+1])

		if zoneid >= int(c.NumZones) {
			// log error?
			continue
		}

		// the prefix fetch is sorted by key in reverse order
		if c.Zones[zoneid] == nil {
			c.Zones[zoneid] = cluster.NewZoneFromConfig(uint32(zoneid), uint32(nodeid+1), c.NumZones, c.NumShards)
		}

		c.Zones[zoneid].Nodes[nodeid].StringToNode(uint32(zoneid), uint32(nodeid),
			string(ev.Value), TagPrimSecondaryDelimiter, TagShardDelimiter)
	}

	return nil
}
