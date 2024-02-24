package proc

import (
	"github.com/golang/glog"
	proto "junodb_lite/pkg/ac_proto"
)

type OnePhaseRequestAndStats struct {
	nextSSIndex int

	//RequestAndStats
	//successResponses []ResponseWrapper
	//errorResponses   []ResponseWrapper
	//
	//nextSSIndex           uint32
	//mostUpdatedOkResponse *ResponseWrapper
}

type CommitRequestAndStats struct {
	//RequestAndStats
	opMsg proto.OperationalMessage
	//noErrResponse           ResponseWrapper
	ssIndicesOfFailedCommit []uint32
}
type RequestAndStats struct {
	raw                 proto.RawMessage
	isSet               bool
	numSent             uint8 //successfully sent
	numPending          uint8
	numFailToSend       uint8
	numFailToSendBusy   uint8 //how much "FailToSend" (numFailToSend) is because of busy
	numFailToSendNoConn uint8 //how much "FailToSend" (numFailToSend) is because of no available connection
	numIOError          uint8 //IO Read Error
	numTimeout          uint8
	numSuccessResponse  uint8
	numErrorResponse    uint8
	funcIsSuccess       func(proto.OpStatus) bool
}

func (s *OnePhaseRequestAndStats) onSuccess(rc *SSRequestContext) {
	switch rc.state {
	case stSSResponseReceived:
		//t := &s.successResponses[s.numSuccessResponse]
		//t.ssRequest = rc
		//s.onSuccessResponse()
		//if rc.ssResponse == nil {
		//	glog.Error("SSRequestContext.ssResponse is nil")
		//}
		//if s.mostUpdatedOkResponse == nil {
		//	s.mostUpdatedOkResponse = t
		//} else {
		//	if recordMostUpdatedThan(&t.ssRequest.ssRespOpMsg, &s.mostUpdatedOkResponse.ssRequest.ssRespOpMsg) {
		//		s.mostUpdatedOkResponse = t
		//	}
		//}
	case stSSRequestIOError, stSSResponseIOError:
		//s.onIOError()
	case stSSRequestTimeout:
		//s.onTimeout()
	default:
		glog.Warning("Unexpected rc.state: ", rc.state)
	}
}
