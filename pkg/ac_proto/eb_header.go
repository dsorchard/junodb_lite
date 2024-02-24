package proto

type (
	opMsgFlagT uint8

	messageHeaderT struct {
		magic    uint16
		version  uint8
		typeFlag messageTypeFlagT
		msgSize  uint32
		opaque   uint32
	}
)
