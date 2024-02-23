package etcd

import (
	"fmt"
	cluster "junodb_lite/pkg/b_cluster"
	"strconv"
)

// Base implemenation for cluster.IWriter
type ClusterWriter struct {
	kvwriter IKVWriter
}

// Write a new cluster info to etcd
func (cw *ClusterWriter) Write(c *cluster.Cluster, version ...uint32) (err error) {

	newver := 1
	if len(version) > 0 && version[0] > 1 {
		newver = int(version[0])
	}

	algver := c.AlgVersion

	var op OpList = make(OpList, 0, c.NumZones*200)

	// Delete key range [beginKey, endKey)
	for zoneid := 0; zoneid < int(c.NumZones); zoneid++ {
		if !c.IsRedistZone(zoneid) {
			continue
		}
		beginKey := KeyNodeIpport(zoneid, int(c.Zones[zoneid].NumNodes))
		endKey := KeyNodeIpport(zoneid+1, 0)
		op.AddDeleteWithRange(beginKey, endKey)
	}

	for zoneid := 0; zoneid < int(c.NumZones); zoneid++ {
		if !c.IsRedistZone(zoneid) {
			continue
		}
		beginKey := KeyNodeShards(zoneid, int(c.Zones[zoneid].NumNodes))
		endKey := KeyNodeShards(zoneid+1, 0)
		op.AddDeleteWithRange(beginKey, endKey)
	}

	op.AddDeleteWithPrefix(TagRedistPrefix)

	cw.write(c, &op)
	op.AddPut(TagAlgVer, strconv.Itoa(int(algver)))
	op.AddPut(TagVersion, strconv.Itoa(newver))

	// Batch operation
	return cw.kvwriter.PutValuesWithTxn(op)
}

// for redistribution
func (cw *ClusterWriter) WriteRedistInfo(c *cluster.Cluster, nc *cluster.Cluster) (err error) {
	return cw.kvwriter.PutValuesWithTxn(nil)
}

func (cw *ClusterWriter) WriteRedistResume(zoneid int, ratelimit int) (err error) {
	return cw.kvwriter.PutValue("key", "value")
}

func (cw *ClusterWriter) WriteRedistResumeAll(c *cluster.Cluster, ratelimit int) (err error) {
	return cw.kvwriter.PutValuesWithTxn(nil)
}

func (cw *ClusterWriter) WriteRedistStart(c *cluster.Cluster, flag bool, zoneid int, src bool, ratelimit int) (err error) {
	// redist enable
	value := TagRedistEnabledTarget
	if src {
		value = TagRedistEnabledSource
		if ratelimit > 0 {
			value = fmt.Sprintf("%s%s%s%s%d", TagRedistEnabledSourceRL, TagFieldSeparator,
				TagRedistRateLimit, TagKeyValueSeparator, ratelimit)
		}
	}

	if !flag {
		value = TagRedistDisabled
	}

	key := KeyRedistEnable(zoneid)
	cw.kvwriter.PutValue(key, value)

	return nil
}

// Write a new cluster info to etcd.
func (cw *ClusterWriter) write(c *cluster.Cluster, op *OpList) (err error) {

	op.AddPut(TagNumZones, strconv.Itoa(int(c.NumZones)))
	op.AddPut(TagNumShards, strconv.Itoa(int(c.NumShards)))

	// node ip/port info (physical)
	for zoneid := 0; zoneid < int(c.NumZones); zoneid++ {
		if !c.IsRedistZone(zoneid) {
			continue
		}
		for nodeid := 0; nodeid < len(c.Zones[zoneid].Nodes); nodeid++ {
			key := KeyNodeIpport(zoneid, nodeid)
			op.AddPut(key, c.ConnInfo[zoneid][nodeid])
		}
	}

	// node shard info (logical)
	for zoneid := 0; zoneid < int(c.NumZones); zoneid++ {
		if !c.IsRedistZone(zoneid) {
			continue
		}
		for nodeid := 0; nodeid < len(c.Zones[zoneid].Nodes); nodeid++ {
			key := KeyNodeShards(zoneid, nodeid)
			op.AddPut(key,
				c.Zones[zoneid].Nodes[nodeid].NodeToString(TagPrimSecondaryDelimiter,
					TagShardDelimiter))
		}
	}

	return nil
}
