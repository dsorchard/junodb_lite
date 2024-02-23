package cli

import proto "junodb_lite/pkg/ac_proto"

// GetResponse() != nil and GetError() != nil are mutually exclusive
type IResponseContext interface {
	GetResponse() *proto.OperationalMessage
	GetError() error
	GetOpaque() uint32
	SetOpaque(opaque uint32)
}

type RequestContext struct {
	request    *proto.OperationalMessage
	chResponse chan IResponseContext
}
