package proc

import (
	"context"
	proto "junodb_lite/pkg/ac_proto"
	io "junodb_lite/pkg/y_conn_mgr"
)

type (
	ProcessorBase struct {
		ctx           context.Context
		clientRequest proto.OperationalMessage
		//repRequest         proto.OperationalMessage
		requestContext io.IRequestContext
		chSSResponse   chan io.IResponseContext
		//ssGroup        SSGroup
		shardId   uint16
		requestID string

		numSSRequestSent      int
		numSSResponseReceived int
		numSSResponseIOError  int

		ssRequestContexts []SSRequestContext

		pendingResponses     []*SSRequestContext
		pendingResponseQueue []*SSRequestContext
		//responseTimer        *util.TimerWrapper
		hasRepliedClient bool

		self IRequestProcessor
	}
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
