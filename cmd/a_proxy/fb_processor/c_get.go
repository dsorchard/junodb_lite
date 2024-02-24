package proc

import (
	proto "junodb_lite/pkg/ac_proto"
	"junodb_lite/pkg/y_conn_mgr"
)

var _ IOnePhaseProcessor = (*GetProcessor)(nil)

type GetProcessor struct {
	OnePhaseProcessor

	repair               RequestAndStats
	numNoKey             int
	numTTLExtendFailures int
}

func (p *GetProcessor) Init() {
	//TODO implement me
	panic("implement me")
}

func (p *GetProcessor) sendInitRequests() {
	//TODO implement me
	panic("implement me")
}

func (p *GetProcessor) OnResponseReceived(rc *SSRequestContext) {
	if rc.opCode == proto.OpCodeRead {
		switch rc.ssResponseOpStatus {
		case proto.OpStatusNoError:
			p.onSuccess(rc)
			return
		//case proto.OpStatusNoKey:
		//	p.onNoKey(rc)
		//case proto.OpStatusBadParam:
		//	p.onFailure(rc)
		case proto.OpStatusKeyMarkedDelete:
			p.onSuccess(rc)
		//case proto.OpStatusSSReadTTLExtendErr:
		//	p.onFailToExtendTTL(rc)
		default:
			//glog.Infof("unexpected response. %s %s", rc.opCode.String(),
			//	rc.ssResponseOpStatus.String())
			//p.onFailure(rc)
		}
	}
}

func (p *GetProcessor) OnSSTimeout(st *SSRequestContext) {
	//TODO implement me
	panic("implement me")
}

func (p *GetProcessor) OnSSIOError(st *SSRequestContext) {
	//TODO implement me
	panic("implement me")
}

func (p *GetProcessor) Process(reqCtx io.IRequestContext) bool {
	//TODO implement me
	panic("implement me")
}

func (p *GetProcessor) setInitSSRequest() bool {
	//TODO implement me
	panic("implement me")
}

func (p *GetProcessor) validateSSResponse(st *SSRequestContext) bool {
	//TODO implement me
	panic("implement me")
}

func (p *GetProcessor) needApplyUDF() bool {
	//TODO implement me
	panic("implement me")
}

func (p *GetProcessor) applyUDF(opmsg *proto.OperationalMessage) {
	//TODO implement me
	panic("implement me")
}

func (p *GetProcessor) succeeded() bool {
	//TODO implement me
	panic("implement me")
}

func (p *GetProcessor) failed() bool {
	//TODO implement me
	panic("implement me")
}

func (p *GetProcessor) sendRequest() {
	//TODO implement me
	panic("implement me")
}

func (p *GetProcessor) onSuccess(rc *SSRequestContext) {
	p.OnePhaseProcessor.onSuccess(rc)

	if p.succeeded() {
		//p.replyToClientAndRepair()
	}
}

func (p *GetProcessor) errorResponseOpStatus() proto.OpStatus {
	//TODO implement me
	panic("implement me")
}

func NewGetProcessor() *GetProcessor {
	p := &GetProcessor{
		OnePhaseProcessor: OnePhaseProcessor{
			ssRequestOpCode: proto.OpCodeRead,
		},
	} //proto.OpCodeGet
	p.self = p
	return p
}
