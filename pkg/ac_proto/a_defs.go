package proto

type (
	OpCode           uint8
	OpStatus         uint8
	messageTypeFlagT uint8
	shardIdOrStatusT [2]uint8
)

func (op OpCode) String() string {
	//return opCodeNameMap[op]
	return ""
}

const (
	OpCodeNop         = OpCode(0)
	OpCodeCreate      = OpCode(1)
	OpCodeGet         = OpCode(2)
	OpCodeUpdate      = OpCode(3)
	OpCodeSet         = OpCode(4)
	OpCodeDestroy     = OpCode(5)
	OpCodeUDFGet      = OpCode(6)
	OpCodeUDFSet      = OpCode(7)
	OpCodeLastProxyOp = OpCode(8) // add proxy op before this

	OpCodePrepareCreate = OpCode(0x81)
	OpCodeRead          = OpCode(0x82)
	OpCodePrepareUpdate = OpCode(0x83)
	OpCodePrepareSet    = OpCode(0x84)
	OpCodePrepareDelete = OpCode(0x85)
	OpCodeDelete        = OpCode(0x86)

	OpCodeCommit     = OpCode(0xC1)
	OpCodeAbort      = OpCode(0xC2)
	OpCodeRepair     = OpCode(0xC3)
	OpCodeMarkDelete = OpCode(0xC4)

	OpCodeClone        = OpCode(0xE1)
	OpCodeVerHandshake = OpCode(0xE2)

	OpCodeMockGetExtendTTL = OpCode(0xFD)
	OpCodeMockSetParam     = OpCode(0xFE)
	OpCodeMockReSet        = OpCode(0xFF)
)
