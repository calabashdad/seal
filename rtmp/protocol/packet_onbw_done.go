package protocol

type OnBwDonePacket struct {
}

func (pkt *OnBwDonePacket) Decode(data []uint8) (err error) {
	return
}
func (pkt *OnBwDonePacket) Encode() (data []uint8) {
	return
}
func (pkt *OnBwDonePacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}
func (pkt *OnBwDonePacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverConnection
}
