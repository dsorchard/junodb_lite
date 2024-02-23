package cluster

import "errors"

type IReader interface {
	Read(c *Cluster) (version uint32, err error)               // for proxy
	ReadWithRedistInfo(c *Cluster) (version uint32, err error) // for storage server
	ReadWithRedistNodeShards(c *Cluster) (err error)           // for cluster manager
}

func (c *Cluster) Read(r IReader) (version uint32, err error) {
	if r == nil {
		return 0, errors.New("nil cluster reader")
	}
	return r.Read(c)
}

func (c *Cluster) ReadWithRedistInfo(r IReader) (version uint32, err error) {
	if r == nil {
		return 0, errors.New("nil cluster reader")
	}
	return r.ReadWithRedistInfo(c)
}

func (c *Cluster) ReadWithRedistNodeShards(r IReader) (err error) {
	if r == nil {
		return errors.New("nil cluster reader")
	}
	return r.ReadWithRedistNodeShards(c)
}

type IWriter interface {
	Write(c *Cluster, version ...uint32) (err error)
	WriteRedistInfo(c *Cluster, nc *Cluster) (err error)                                     // write redistribution info
	WriteRedistStart(c *Cluster, flag bool, zoneid int, src bool, ratelimit int) (err error) // write redistibution start/stop
	WriteRedistAbort(c *Cluster) (err error)
	WriteRedistResume(zoneid int, ratelimit int) (err error)
}

func (c *Cluster) Write(w IWriter, version ...uint32) (err error) {
	if w == nil {
		return errors.New("nil cluster writer")
	}
	if len(version) > 0 {
		return w.Write(c, version[0])
	} else {
		return w.Write(c)
	}
}

func (c *Cluster) WriteRedistInfo(w IWriter, nc *Cluster) (err error) {
	if w == nil {
		return errors.New("nil cluster writer")
	}

	return w.WriteRedistInfo(c, nc)
}

func (c *Cluster) WriteRedistStart(w IWriter, flag bool, zoneid int, src bool, ratelimit int) (err error) {
	if w == nil {
		return errors.New("nil cluster writer")
	}

	return w.WriteRedistStart(c, flag, zoneid, src, ratelimit)
}

func (c *Cluster) WriteRedistAbort(w IWriter) (err error) {
	if w == nil {
		return errors.New("nil cluster writer")
	}

	return w.WriteRedistAbort(c)
}
