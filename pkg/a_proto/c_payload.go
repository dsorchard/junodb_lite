package proto

type (
	PayloadType uint8

	Payload struct {
		tag  PayloadType
		data []byte
	}
)

func (p *Payload) GetLength() uint32 {
	szPayload := len(p.data)
	if szPayload == 0 {
		return 0
	}
	return uint32(1 + szPayload)
}

func (p *Payload) SetWithClearValue(value []byte) {

}
