package proto

func GetOpCode(wmsg *RawMessage) (opCode OpCode, err error) {
	//if wmsg.getMsgType() != kOperationalMessageType || len(wmsg.body) < kOpMsgSubHeaderSize {
	//	err = &ProtocolError{"not Operational Message"}
	//	return
	//}
	opCode = OpCode(wmsg.body[0])
	return
}
