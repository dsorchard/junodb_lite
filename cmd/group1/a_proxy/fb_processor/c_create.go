package proc

import (
	proto "junodb_lite/pkg/ac_proto"
	"junodb_lite/pkg/y_conn_mgr"
)

type CreateProcessor struct {
	TwoPhaseProcessor
	numDupKey    uint8
	numInserting uint8
}

func (p *CreateProcessor) Init() {
	//TODO implement me
	panic("implement me")
}

func (p *CreateProcessor) sendInitRequests() {
	//TODO implement me
	panic("implement me")
}

func (p *CreateProcessor) OnResponseReceived(st *SSRequestContext) {
	//TODO implement me
	panic("implement me")
}

func (p *CreateProcessor) OnSSTimeout(st *SSRequestContext) {
	//TODO implement me
	panic("implement me")
}

func (p *CreateProcessor) OnSSIOError(st *SSRequestContext) {
	//TODO implement me
	panic("implement me")
}

func (p *CreateProcessor) Process(reqCtx io.IRequestContext) bool {
	//TODO implement me
	panic("implement me")
}

func (p *CreateProcessor) setInitSSRequest() bool {
	//TODO implement me
	panic("implement me")
}

func (p *CreateProcessor) validateSSResponse(st *SSRequestContext) bool {
	//TODO implement me
	panic("implement me")
}

func (p *CreateProcessor) needApplyUDF() bool {
	//TODO implement me
	panic("implement me")
}

func (p *CreateProcessor) applyUDF(opmsg *proto.OperationalMessage) {
	//TODO implement me
	panic("implement me")
}

func NewCreateProcessor() *CreateProcessor {
	p := &CreateProcessor{
		TwoPhaseProcessor: TwoPhaseProcessor{
			prepareOpCode: proto.OpCodePrepareCreate,
		},
	}
	p.self = p
	return p
}
