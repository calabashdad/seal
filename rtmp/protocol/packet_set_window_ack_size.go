package protocol

type SetWindowAckSizePacket struct {
}

func (pkt *SetWindowAckSizePacket) Decode(data []uint8) (err error) {
	return
}
func (pkt *SetWindowAckSizePacket) Encode() (data []uint8) {
	return
}
func (pkt *SetWindowAckSizePacket) GetMessageType() uint8 {
	return RTMP_MSG_WindowAcknowledgementSize
}
func (pkt *SetWindowAckSizePacket) GetPreferCsId() uint32 {
	return RTMP_CID_ProtocolControl
}
