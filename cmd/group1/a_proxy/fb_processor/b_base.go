package proc

import (
	proto "junodb_lite/pkg/ac_proto"
	io "junodb_lite/pkg/y_conn_mgr"
	"time"
)

type (

	// virtual functions, to trigger polymorphism, call via p.self pointer.
	IRequestProcessor interface {
		Init()
		sendInitRequests()

		OnResponseReceived(st *SSRequestContext)
		OnSSTimeout(st *SSRequestContext)
		OnSSIOError(st *SSRequestContext)

		//Return true if it has completed all the SS requests and can be cached
		Process(reqCtx io.IRequestContext) bool

		setInitSSRequest() bool
		validateSSResponse(st *SSRequestContext) bool

		needApplyUDF() bool
		applyUDF(opmsg *proto.OperationalMessage)
	}

	SSRequestContext struct {
		timeToExpire     time.Time
		timeReqSent      time.Time
		timeRespReceived time.Time // when state is changed from stSSRequestSent to others

		ssRequest          io.IRequestContext
		ssResponse         io.IResponseContext
		ssRespOpMsg        proto.OperationalMessage
		opCode             proto.OpCode
		ssResponseOpStatus proto.OpStatus
		ssIndex            uint32
		//state              ssReqContextState
	}
)
