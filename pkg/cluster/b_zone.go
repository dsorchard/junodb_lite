package cluster

// Node class represent a logic node
type Zone struct {
	Zoneid   uint32
	NumNodes uint32
	Nodes    []Node
}
