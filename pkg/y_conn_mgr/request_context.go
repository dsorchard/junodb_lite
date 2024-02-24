package io

import (
	"context"
	"io"
	proto "junodb_lite/pkg/ac_proto"
	util "junodb_lite/pkg/y_util"
	"time"
)

type (
	IRequestContext interface {
		util.QueItem
		GetMessage() *proto.RawMessage
		GetCtx() context.Context
		Cancel()
		Read(r io.Reader) (n int, err error)
		WriteWithOpaque(opaque uint32, w io.Writer) (n int, err error)
		Reply(resp IResponseContext)
		OnComplete()
		GetReceiveTime() time.Time
		SetTimeout(parent context.Context, duration time.Duration)
	}
	IResponseContext interface {
		GetStatus() uint32
		GetMessage() *proto.RawMessage
		GetMsgSize() uint32
		OnComplete()
		Read(r io.Reader) (n int, err error)
		Write(w io.Writer) (n int, err error)
	}
)
