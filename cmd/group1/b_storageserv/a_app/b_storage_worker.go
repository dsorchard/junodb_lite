package app

import util "junodb_lite/pkg/y_util"

type Worker struct {
	CmdStorageCommon
	optWorkerId        uint
	optListenAddresses util.StringListFlags
	optIsChild         bool
	optHttpMonAddr     string
	optZoneId          uint
	optMachineIndex    uint
	optLRUCacheSize    uint
}

func (c *Worker) GetName() string {
	//TODO implement me
	panic("implement me")
}

func (c *Worker) GetDesc() string {
	//TODO implement me
	panic("implement me")
}

func (c *Worker) GetSynopsis() string {
	//TODO implement me
	panic("implement me")
}

func (c *Worker) GetDetails() string {
	//TODO implement me
	panic("implement me")
}

func (c *Worker) GetOptionDesc() string {
	//TODO implement me
	panic("implement me")
}

func (c *Worker) GetExample() string {
	//TODO implement me
	panic("implement me")
}

func (c *Worker) AddExample(cmdExample string, desc string) {
	//TODO implement me
	panic("implement me")
}

func (c *Worker) AddDetails(txt string) {
	//TODO implement me
	panic("implement me")
}

func (c *Worker) Exec() {
	//TODO implement me
	panic("implement me")
}

func (c *Worker) PrintUsage() {
	//TODO implement me
	panic("implement me")
}

func (c *Worker) Init(name string, desc string) {
	//c.CmdStorageCommon.Init(name, desc)
	//c.UintOption(&c.optWorkerId, "id|worker-id", 0, "specify the ID of the worker")
	//c.ValueOption(&c.optListenAddresses, "listen", "specify listening address. Override Listener in config file")
	//c.BoolOption(&c.optIsChild, "child", false, "specify if the worker was started by a parent process")
	//c.StringOption(&c.optHttpMonAddr, "mon-addr|monitoring-address", "", "specify the http monitoring address. \n\toverride HttpMonAddr in config file")
	//c.UintOption(&c.optZoneId, "zone-id", 0, "specify zone id")
	//c.UintOption(&c.optMachineIndex, "machine-index", 0, "specify machine index")
	//c.UintOption(&c.optLRUCacheSize, "lru-cache-mb", 0, "specify lru cache size")
}
