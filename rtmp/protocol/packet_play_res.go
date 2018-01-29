package protocol

type PlayResPacket struct {
}

func (pkt *PlayResPacket) Decode(data []uint8) (err error) {
	return
}
func (pkt *PlayResPacket) Encode() (data []uint8) {
	return
}
func (pkt *PlayResPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}
func (pkt *PlayResPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverStream
}
