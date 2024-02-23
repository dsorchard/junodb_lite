package proto

type OperationalMessage struct {
	opCode OpCode
	//flags           opMsgFlagT
	shardIdOrStatus shardIdOrStatusT
	typeFlag        messageTypeFlagT
	opaque          uint32

	namespace []byte
	key       []byte
	payload   Payload

	//timeToLive           timeToLiveT
	//version              versionT
	//creationTime         creationTimeT
	//expirationTime       expirationTimeT
	//lastModificationTime lastModificationTimeT
	//sourceInfo           sourceInfoT
	//requestID            requestIdT
	//originatorRequestID  originatorT
	//correlationID        correlationIdT
	//requestHandlingTime  requestHandlingTimeT
	//udfName              udfNameT
}

func (m *OperationalMessage) SetCorrelationID(id []byte) {
	//m.correlationID.set(id)
}

func (m *OperationalMessage) GetPayload() *Payload {
	return &m.payload
}

func (m *OperationalMessage) SetRequest(op OpCode, key []byte, bytes []byte, p *Payload, ttl uint32) {

}

func (m *OperationalMessage) SetNewRequestID() {

}

func (m *OperationalMessage) SetUDFName(fname []byte) {

}

func (m *OperationalMessage) GetOpCode() OpCode {
	return m.opCode
}

func (m *OperationalMessage) GetOpStatus() OpStatus {
	//if m.typeFlag.isResponse() {
	//	return OpStatus(m.shardIdOrStatus[1])
	//} else {
	//	///log
	//	return OpStatusNoError
	//}
	return 0
}

func (m *OperationalMessage) GetOpCodeText() any {
	return nil
}
