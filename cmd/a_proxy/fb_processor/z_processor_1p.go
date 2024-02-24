package proc

import (
	proto "junodb_lite/pkg/ac_proto"
)

type (
	IOnePhaseProcessor interface {
		IRequestProcessor
		succeeded() bool
		failed() bool
		sendRequest()
		onSuccess(rc *SSRequestContext)
		errorResponseOpStatus() proto.OpStatus
	}
	OnePhaseProcessor struct {
		ProcessorBase
		request         OnePhaseRequestAndStats
		ssRequestOpCode proto.OpCode
	}
	OnePhaseRequestAndStats struct {
		RequestAndStats
		//successResponses []ResponseWrapper
		//errorResponses   []ResponseWrapper

		nextSSIndex uint32
		//mostUpdatedOkResponse *ResponseWrapper
	}
)

func (p *OnePhaseProcessor) onSuccess(rc *SSRequestContext) {
	p.request.onSuccess(rc)
}

func (p *OnePhaseProcessor) sendRequest() {
	if int(p.request.nextSSIndex) < p.ssGroup.numAvailableSSs {
		p.send(&p.request.RequestAndStats, p.request.nextSSIndex)
		p.request.nextSSIndex++
	}
}
