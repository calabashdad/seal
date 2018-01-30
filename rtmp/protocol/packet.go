package protocol

//Packet is the interface of the packet.
type Packet interface {
	Decode([]uint8) error
	Encode() []uint8
	GetMessageType() uint8
	GetPreferCsId() uint32
}
