package proc

import (
	proto "junodb_lite/pkg/ac_proto"
)

type (
	TwoPhaseProcessor struct {
		ProcessorBase
		prepareOpCode proto.OpCode
		prepare       OnePhaseRequestAndStats
		state         twoPhaseProcessorState

		numBadRequestID int
		commit          CommitRequestAndStats

		abort  RequestAndStats
		repair RequestAndStats
	}
)

type IOnePhaseProcessor interface {
	IRequestProcessor
	succeeded() bool
	failed() bool
	sendRequest()
	onSuccess(rc *SSRequestContext)
	errorResponseOpStatus() proto.OpStatus
}
type OnePhaseProcessor struct {
	ProcessorBase
	request         OnePhaseRequestAndStats
	ssRequestOpCode proto.OpCode
}

func (p *OnePhaseProcessor) onSuccess(rc *SSRequestContext) {
	p.request.onSuccess(rc)
}
