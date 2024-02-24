package proc

import proto "junodb_lite/pkg/ac_proto"

type (
	ITwoPhaseProcessor interface {
		IRequestProcessor
		onPrepareSuccess(rc *SSRequestContext)
		prepareSucceeded() bool
		prepareFailed() bool
		commitSucceeded() bool
		sendCommit(ssIndex uint32)
		sendAbort(ssIndex uint32)
		sendRepair(ssIndex uint32)
		sendPrepareRequest()
		abortSucceededPrepares()
		//Returns the OpStatus to Client
		errorPrepareResponseOpStatus() proto.OpStatus
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

	twoPhaseProcessorState uint8
)

const (
	stTwoPhaseProcInit twoPhaseProcessorState = iota
	stTwoPhaseProcPrepare
	stTwoPhaseProcCommit
	stTwoPhaseProcAbort
)
