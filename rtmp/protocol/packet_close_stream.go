package protocol

type CloseStreamPacket struct {
}

func (pkt *CloseStreamPacket) Decode([]uint8) (err error) {
	return
}
func (pkt *CloseStreamPacket) Encode() (b []uint8) {
	return
}
func (pkt *CloseStreamPacket) GetMessageType() uint8 {
	//no method for this pakcet
	return 0
}
func (pkt *CloseStreamPacket) GetPreferCsId() uint32 {
	//no method for this pakcet
	return 0
}
