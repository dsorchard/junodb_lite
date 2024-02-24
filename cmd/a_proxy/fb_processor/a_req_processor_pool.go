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
func (p *ReqProcessorPool) GetProcessor() IRequestProcessor {

	//// reached absolute max, should reject or queue request
	//if p.GetCount() >= p.maxCount {
	//	return nil
	//}
	//
	//if p.curCount != nil {
	//	p.curCount.Add(1)
	//}
	return p.procPool.Get().(IRequestProcessor)
}
func (p *ReqProcessorPool) PutProcessor(proc IRequestProcessor) {
	proc.Init()
	p.procPool.Put(proc)
	//if p.curCount != nil {
	//	p.curCount.Add(-1)
	//}
}

func (p *ReqProcessorPool) DecreaseCount() {
	//if p.curCount != nil {
	//	p.curCount.Add(-1)
	//}
}
