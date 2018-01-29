package protocol

type PlayPacket struct {
}

func (pkt *PlayPacket) Decode([]uint8) (err error) {
	return
}
func (pkt *PlayPacket) Encode() (b []uint8) {
	return
}
func (pkt *PlayPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}
func (pkt *PlayPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverStream
}
