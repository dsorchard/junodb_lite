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

type StatsConfig struct {
	TimeoutStatsEnabled    bool
	RespTimeStatsEnabled   bool
	MarkdownThreashold     uint32
	MarkdownExpirationBase uint32
	EMARespTimeWindowSize  uint32
	TimeoutWindowSize      uint32
	TimeoutWindowUint      uint32
}
