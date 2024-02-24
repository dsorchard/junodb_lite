package proc

import (
	"context"
	"github.com/golang/glog"
	proto "junodb_lite/pkg/ac_proto"
	cluster "junodb_lite/pkg/b_cluster"
	shard "junodb_lite/pkg/b_shard"
	io "junodb_lite/pkg/y_conn_mgr"
	"time"
)

const (
	stSSRequestInit ssReqContextState = iota
	stSSRequestSent
	stSSResponseReceived
	stSSRequestIOError
	stSSResponseIOError
	stSSRequestTimeout
	stRequestTimeout
	stRequestCancelled
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

	ProcessorBase struct {
		ctx           context.Context
		clientRequest proto.OperationalMessage
		//repRequest         proto.OperationalMessage
		requestContext io.IRequestContext
		chSSResponse   chan io.IResponseContext
		ssGroup        SSGroup
		shardId        uint16
		requestID      string

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
	SSGroup struct {
		processors      []*cluster.OutboundSSProcessor
		procIndices     []int
		numAvailableSSs int
		numBrokenSSs    int
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
		state              ssReqContextState

		//state              ssReqContextState
	}
	ssReqContextState uint8
)

func (g *SSGroup) getProcessors(key []byte) (shardId shard.ID, ok bool) {
	shardId, g.numAvailableSSs = cluster.GetShardMgr().GetSSProcessors(key, 1, g.processors, g.procIndices)
	//g.numBrokenSSs = confNumZones - g.numAvailableSSs
	//ok = g.numAvailableSSs >= confNumWrites
	return
}
func (p *ProcessorBase) Process(request io.IRequestContext) bool {

	p.ctx = request.GetCtx()
	p.requestContext = request
	p.clientRequest = proto.OperationalMessage{}

	//if err := p.clientRequest.Decode(request.GetMessage()); err != nil {
	//	glog.Error("Failed to decode inbound request: ", err)
	//	p.replyStatusToClient(proto.OpStatusBadMsg)
	//	p.OnComplete()
	//	return true
	//}

	//p.requestID = p.clientRequest.GetRequestIDString()
	shardId, ok := p.ssGroup.getProcessors(p.clientRequest.GetKey())

	if !ok {
		//p.replyStatusToClient(proto.OpStatusNoStorageServer)
		//glog.Warning("Cannot get channels from Cluster Manager")
		//p.OnComplete()
		//return true
	}

	p.shardId = shardId.Uint16()

	if err := proto.SetShardId(p.requestContext.GetMessage(), p.shardId); err != nil {
		//p.replyStatusToClient(proto.OpStatusInternal) //shouldn't happen.
		glog.Error("fail to set shardId: ", err)
		return true
	}

	p.self.sendInitRequests()
	done := false

loop:
	for p.isDone() == false {
		select {
		case <-p.ctx.Done():
			if done == false {
				done = true
				if p.ctx.Err() == context.DeadlineExceeded {
					p.OnRequestTimeout()
				} else {
					p.OnCancelled()
				}
			}
			break loop
		//case t := <-p.chSSTimeout():
		//	p.handleSSTimeout(t)
		case respFromSS := <-p.chSSResponse:
			p.onResponseReceived(respFromSS)
		}
	}
	if p.isDone() {
		//p.OnComplete()
		return true
	}
	return false
}

func (p *ProcessorBase) isDone() bool {
	return (p.numSSRequestSent == p.numSSResponseReceived)
}

func (p *ProcessorBase) OnRequestTimeout() {

}

func (p *ProcessorBase) OnCancelled() {

}

func (p *ProcessorBase) onResponseReceived(resp io.IResponseContext) {
	st := p.preprocessAndValidateResponse(resp)
	if st != nil {
		//st.timeRespReceived = time.Now()
		p.self.OnResponseReceived(st)
	} else {
		//io.ReleaseOutboundResponse(resp)
	}
}
func (p *ProcessorBase) preprocessAndValidateResponse(resp io.IResponseContext) (st *SSRequestContext) {
	return
}
