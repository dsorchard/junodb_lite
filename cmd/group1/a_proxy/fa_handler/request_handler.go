package handler

import (
	"github.com/golang/glog"
	config "junodb_lite/cmd/group1/a_proxy/b_config"
	proc "junodb_lite/cmd/group1/a_proxy/fb_processor"
	proto "junodb_lite/pkg/ac_proto"
	service "junodb_lite/pkg/g_service_mgr"
	"junodb_lite/pkg/y_conn_mgr"
	"os"
)

type RequestHandler struct {
	procPools []*proc.ReqProcessorPool
}

func (rh *RequestHandler) Init() {
	//TODO implement me
	panic("implement me")
}

func (rh *RequestHandler) GetReqCtxCreator() io.InboundRequestContextCreator {
	//TODO implement me
	panic("implement me")
}

func (rh *RequestHandler) Process(reqCtx io.IRequestContext) error {
	op, err := proto.GetOpCode(nil)
	if err != nil {
		glog.Error("Cannot get Opcode: ", err)
		return err
	}
	processor := rh.GetProcessor(op)
	if processor == nil {
		glog.Error("Cannot get processor Opcode: ", op)
		return nil
	}

	if processor.Process(reqCtx) {
		rh.PutProcessor(op, processor)
	} else {
		/// not put back to the processor pool, as the channel(s) might still be referenced by OutboundProcessor.
		rh.decreaseProcPoolCount(op)
	}
	return nil
}

func (rh *RequestHandler) GetProcessor(op proto.OpCode) proc.IRequestProcessor {
	procPool := rh.procPools[op]
	if procPool == nil {
		return nil
	}
	return procPool.GetProcessor()
}
func (rh *RequestHandler) PutProcessor(op proto.OpCode, p proc.IRequestProcessor) {
	procPool := rh.procPools[op]
	if procPool != nil {
		procPool.PutProcessor(p)
	}
}

func (rh *RequestHandler) decreaseProcPoolCount(op proto.OpCode) {
	procPool := rh.procPools[op]
	if procPool != nil {
		procPool.DecreaseCount()
	}
}

func (rh *RequestHandler) Finish() {
	//TODO implement me
	panic("implement me")
}

func (rh *RequestHandler) OnKeepAlive(connector *io.Connector, reqCtx io.IRequestContext) error {
	//TODO implement me
	panic("implement me")
}

func NewRequestHandler() *RequestHandler {
	rh := &RequestHandler{}
	return rh
}

func NewProxyServiceWithListenFd(conf *config.Config, limiter service.ILimiter, fds ...*os.File) *service.Service {
	s := service.NewWithLimiterAndListenFd(conf.Config, NewRequestHandler(), limiter, fds...)
	return s
}

func NewProxyService(conf *config.Config) *service.Service {
	s, _ := service.NewService(conf.Config, NewRequestHandler())

	//stats.SetListeners(s.GetListeners())
	return s
}
