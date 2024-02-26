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
		z.addOneNode(nil)
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

func (z *Zone) Log() {
	glog.Infof("zoneid=%d, numNodes=%d", z.Zoneid, z.NumNodes)
	for i := uint32(0); i < z.NumNodes; i++ {
		z.Nodes[i].Log()
	}
}

// helper
type indexElem struct {
	nodeid int
	weight int
}

func (z *Zone) addOneNode(filter *Filter) (err error) {
	var target Node
	target.InitNode(z.Zoneid, uint32(len(z.Nodes)))

	var curr int = 0
	var numNodes = len(z.Nodes)
	var reorder []indexElem = make([]indexElem, numNodes)
	for k := 0; k < numNodes; k++ {
		reorder[k].nodeid = k
		//reorder[k].weight = z.Nodes[k].primaryLength()
	}

	//if !IsNewMappingAlg() {
	//	sort.Sort(byWeight(reorder))
	//}

	for {
		nodeid := reorder[curr].nodeid

		if IsNewMappingAlg() {
			max := 0
			nodeid = 0
			for j := 0; j < numNodes; j++ {
				if z.Nodes[j].primaryLength() >= max {
					max = z.Nodes[j].primaryLength()
					nodeid = j
				}
			}
		}

		source := z.Nodes[nodeid]

		// Check exit condition
		if target.primaryLength() >= source.primaryLength()-1 {
			break
		}

		shardid := z.Nodes[nodeid].deleteFromPrimary(filter, target)
		target.appendToPrimary(shardid)

		if IsNewMappingAlg() {
			continue
		}

		next := curr + 1

		if next >= numNodes {
			curr = 0
			continue
		}

		// Rewind if the current node has more shards than next.
		nextNode := reorder[next].nodeid
		if z.Nodes[nodeid].primaryLength() >= z.Nodes[nextNode].primaryLength() {
			curr = 0 // rewind
		} else { // move to the next
			curr = next
		}
	}

	if IsNewMappingAlg() {
		//sort.Sort(byValue(target.PrimaryShards))
		z.Nodes = append(z.Nodes, target)

		return nil
	}

	// TODO reuse code
	for k := 0; k < numNodes; k++ {
		reorder[k].nodeid = k
		reorder[k].weight = z.Nodes[k].secondaryLength()
	}

	//sort.Sort(byWeight(reorder))

	curr = 0
	for {
		nodeid := reorder[curr].nodeid
		source := z.Nodes[nodeid]

		// Check exit condition
		if target.secondaryLength() >= source.secondaryLength()-1 {
			break
		}

		shardid := z.Nodes[nodeid].deleteFromSecondary(filter, target)
		target.appendToSecondary(shardid)

		next := curr + 1
		if next >= numNodes {
			curr = 0
			continue
		}

		// Rewind if the current node has more shards than next.
		nextNode := reorder[next].nodeid
		if z.Nodes[nodeid].secondaryLength() >= z.Nodes[nextNode].secondaryLength() {
			curr = 0 // rewind
		} else { // move to the next
			curr = next
		}
	}

	z.Nodes = append(z.Nodes, target)
	return nil
}

func (z *Zone) removeOneNode() {

	var curr int = 0
	var last = len(z.Nodes) - 1
	if last <= 0 {
		return // done
	}

	var reorder []indexElem = make([]indexElem, last)
	for k := 0; k < last; k++ {
		reorder[k].nodeid = k
		reorder[k].weight = -z.Nodes[k].primaryLength()
	}
	//sort.Sort(byWeight(reorder))

	for {
		// Target nodeid
		nodeid := reorder[curr].nodeid

		// Check exit condition
		if z.Nodes[last].primaryLength() == 0 {
			break
		}

		shardid := z.Nodes[last].deleteFromPrimary(nil, z.Nodes[nodeid])
		z.Nodes[nodeid].appendToPrimary(shardid)

		next := curr + 1
		if next >= last {
			curr = 0
			continue
		}

		// Rewind if the current node has fewer shards than next.
		nextNode := reorder[next].nodeid
		if z.Nodes[nodeid].primaryLength() <= z.Nodes[nextNode].primaryLength() {
			curr = 0 // rewind
		} else { // move to the next
			curr = next
		}
	}

	for k := 0; k < last; k++ {
		reorder[k].nodeid = k
		reorder[k].weight = -z.Nodes[k].secondaryLength()
	}

	//sort.Sort(byWeight(reorder))

	curr = 0
	for {
		// Target nodeid
		nodeid := reorder[curr].nodeid

		// Check exit condition
		if z.Nodes[last].secondaryLength() == 0 {
			break
		}

		shardid := z.Nodes[last].deleteFromSecondary(nil, z.Nodes[nodeid])
		z.Nodes[nodeid].appendToSecondary(shardid)

		next := curr + 1
		if next >= last {
			curr = 0
			continue
		}

		// Rewind if the current node has fewer shards than next.
		nextNode := reorder[next].nodeid
		if z.Nodes[nodeid].secondaryLength() <= z.Nodes[nextNode].secondaryLength() {
			curr = 0 // rewind
		} else { // move to the next
			curr = next
		}
	}

	// Remove last node
	z.Nodes = z.Nodes[:last]
}
