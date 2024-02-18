package cluster

// Node class represent a logic node
type Node struct {
	Zoneid          uint32
	Nodeid          uint32
	PrimaryShards   []uint32
	SecondaryShards []uint32
}
