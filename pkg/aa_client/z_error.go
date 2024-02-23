package client

import proto "junodb_lite/pkg/ac_proto"

var errorMapping map[proto.OpStatus]error

var (
	ErrNoKey              error
	ErrUniqueKeyViolation error
	ErrBadParam           error
	ErrConditionViolation error

	ErrBadMsg           error
	ErrNoStorage        error
	ErrRecordLocked     error
	ErrTTLExtendFailure error
	ErrBusy             error

	ErrWriteFailure   error
	ErrInternal       error
	ErrOpNotSupported error
)
