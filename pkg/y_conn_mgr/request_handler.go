package io

type IRequestHandler interface {
	Init()
	GetReqCtxCreator() InboundRequestContextCreator

	Process(reqCtx IRequestContext) error
	Finish()

	OnKeepAlive(connector *Connector, reqCtx IRequestContext) error
}

type (
	InboundRequestContextCreator func(magic []byte, c *Connector) (ctx IRequestContext, err error)
)
