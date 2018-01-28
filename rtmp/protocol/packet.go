package protocol

type Packet interface {
	Encode() []uint8
	GetMessageType() uint8
	GetPreferCsId() uint32
}
