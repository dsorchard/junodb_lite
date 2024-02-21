package cluster

import (
	"github.com/golang/glog"
)

// Node class represent a logic node
type Zone struct {
	Zoneid   uint32
	NumNodes uint32
	Nodes    []Node
}

var newMappingAlg = false

func IsNewMappingAlg() bool {
	return newMappingAlg
}

func NewZoneFromConfig(zoneid uint32, numNodes uint32, numZones uint32, numShards uint32) *Zone {
	zone := Zone{
		Zoneid:   zoneid,
		NumNodes: numNodes,
		Nodes:    make([]Node, 1, numNodes),
	}

	// Populate Nodes
	zone.initShardsAsssignment(numZones, numShards)
	return &zone
}

func (z *Zone) initShardsAsssignment(numZones uint32, numShards uint32) {

	z.Nodes[0].InitNode(z.Zoneid, 0)
	z.Nodes[0].initShards(z.Zoneid, numZones, numShards)

	for k := uint32(1); k < z.NumNodes; k++ {
		z.addOneNode()
	}
}

func IsPrimary(shardid uint32, zoneid uint32, numZones uint32) bool {

	if IsNewMappingAlg() {
		return true
	}
	m := shardid % numZones

	if m >= zoneid && m < zoneid+(numZones-1)/2 {
		return false
	}

	if m+numZones < zoneid+(numZones-1)/2 {
		return false
	}

	return true
}

func SetMappingAlg(algVersion uint32) {
	glog.Infof("algver=%d", algVersion)
	if algVersion < 2 {
		newMappingAlg = false
		return
	}

	newMappingAlg = true
}

func (z *Zone) addOneNode() (err error) {
	return nil
}
