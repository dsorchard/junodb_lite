package cluster

type Config struct {
	AlgVersion uint32
	NumZones   uint32
	NumShards  uint32
	ConnInfo   [][]string

	//SSHosts and SSPorts are used to generate ConnInfo ONLY when ConnInfo not defined
	SSHosts [][]string
	SSPorts []uint16
}
