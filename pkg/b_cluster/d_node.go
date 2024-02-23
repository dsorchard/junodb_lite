package cluster

import (
	"github.com/golang/glog"
	"strconv"
	"strings"
)

// Node class represent a logic node
type Node struct {
	Zoneid          uint32
	Nodeid          uint32
	PrimaryShards   []uint32
	SecondaryShards []uint32
}

func (n *Node) InitNode(zoneid uint32, nodeid uint32) {
	n.Zoneid = zoneid
	n.Nodeid = nodeid
}

// First node only
func (n *Node) initShards(zoneid uint32, numZones uint32, numShards uint32) {

	var primary []uint32 = make([]uint32, 0, numShards)
	var secondary []uint32 = make([]uint32, 0, numShards)

	for k := uint32(0); k < numShards; k++ {

		if IsPrimary(k, zoneid, numZones) {
			primary = append(primary, k)
		} else {
			secondary = append(secondary, k)
		}
	}
	n.allocate(len(primary), len(secondary))
	n.fillPrimary(0, len(primary), primary)
	n.fillSecondary(0, len(secondary), secondary)
}

func (n *Node) StringToNode(zoneid uint32, nodeid uint32, val string,
	priSecDelimiter string, shardDelimiter string) error {
	return nil
}

func (n *Node) Log() {
	glog.Infof("zoneid=%d, nodeid=%d, prim_shards=%#v, second_shards=%#v",
		n.Zoneid, n.Nodeid, n.PrimaryShards, n.SecondaryShards)
}

func (n *Node) NodeToString(priSecDelimiter string, shardDelimiter string) string {

	var shards_str []string = make([]string, 2)

	var list []string = make([]string, 0, len(n.PrimaryShards))
	for _, shardid := range n.PrimaryShards {
		list = append(list, strconv.Itoa(int(shardid)))
	}
	shards_str[0] = strings.Join(list, shardDelimiter)

	list = make([]string, 0, len(n.SecondaryShards))
	for _, shardid := range n.SecondaryShards {
		list = append(list, strconv.Itoa(int(shardid)))
	}
	shards_str[1] = strings.Join(list, shardDelimiter)

	return strings.Join(shards_str, priSecDelimiter)
}

func (n *Node) GetShards() (shards []uint32) {
	shards = make([]uint32, len(n.PrimaryShards)+len(n.SecondaryShards))
	copy(shards, n.PrimaryShards)
	copy(shards[len(n.PrimaryShards):], n.SecondaryShards)
	return
}
