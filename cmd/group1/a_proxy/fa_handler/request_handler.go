package handler

import (
	config "junodb_lite/cmd/group1/a_proxy/b_config"
	service "junodb_lite/pkg/g_service_mgr"
	"junodb_lite/pkg/y_conn_mgr"
	"os"
)

type RequestHandler struct {
	//procPools []*proc.ReqProcessorPool
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
	//TODO implement me
	panic("implement me")
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
