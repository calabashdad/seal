package protocol

type PublishPacket struct {
}

func (pkt *PublishPacket) Decode([]uint8) (err error) {
	return
}
func (pkt *PublishPacket) Encode() (b []uint8) {
	return
}
func (pkt *PublishPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}
func (pkt *PublishPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverStream
}
