package protocol

type OnStatusCallPacket struct {
}

func (pkt *OnStatusCallPacket) Decode([]uint8) (err error) {
	return
}
func (pkt *OnStatusCallPacket) Encode() (b []uint8) {
	return
}
func (pkt *OnStatusCallPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}
func (pkt *OnStatusCallPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverStream
}
