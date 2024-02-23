package proto

import "io"

type (
	IRequestId interface {
		Bytes() []byte
		String() string
		PrettyPrint(w io.Writer)
	}
	RequestId [16]byte
)
