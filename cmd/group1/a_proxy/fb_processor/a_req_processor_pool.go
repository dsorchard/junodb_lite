package proc

import (
	proto "junodb_lite/pkg/ac_proto"
	util "junodb_lite/pkg/y_util"
)

type ReqProcessorPool struct {
	procPool *util.ChanPool
	maxCount int32
}

func NewRequestProcessorPool(chansize int32, maxsize int32, op proto.OpCode) *ReqProcessorPool {

	procPool := util.NewChanPool(int(chansize), func() interface{} {
		var p IRequestProcessor

		switch op {
		case proto.OpCodeCreate:
			p = NewCreateProcessor()
		case proto.OpCodeGet:
			p = NewGetProcessor()
		//case proto.OpCodeUpdate:
		//	p = NewUpdateProcessor()
		//case proto.OpCodeSet:
		//	p = NewSetProcessor()
		default:
			return nil
		}
		p.Init()
		return p
	})

	return &ReqProcessorPool{procPool, maxsize}
}
