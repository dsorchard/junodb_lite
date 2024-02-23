package client

import "io"

type IContext interface {
	GetVersion() uint32
	GetCreationTime() uint32
	GetTimeToLive() uint32
	PrettyPrint(w io.Writer)
}

///TODO check API input arguments

type IClient interface {
	Create(key []byte, value []byte, opts ...IOption) (IContext, error)
	Get(key []byte, opts ...IOption) ([]byte, IContext, error)
	Update(key []byte, value []byte, opts ...IOption) (IContext, error)
	Set(key []byte, value []byte, opts ...IOption) (IContext, error)
	Destroy(key []byte, opts ...IOption) (err error)
	UDFGet(key []byte, fname []byte, params []byte, opts ...IOption) ([]byte, IContext, error)
	UDFSet(key []byte, fname []byte, params []byte, opts ...IOption) (IContext, error)
}
