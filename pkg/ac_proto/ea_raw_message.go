package proto

type RawMessage struct {
	messageHeaderT
	body []byte
	//buf  *util.PPBuffer
	//pool util.BufferPool
}
