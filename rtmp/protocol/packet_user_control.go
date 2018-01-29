package protocol

type UserControlPacket struct {
}

func (pkt *UserControlPacket) Decode(data []uint8) (err error) {
	return
}
func (pkt *UserControlPacket) Encode() (data []uint8) {
	return
}
func (pkt *UserControlPacket) GetMessageType() uint8 {
	return RTMP_MSG_UserControlMessage
}
func (pkt *UserControlPacket) GetPreferCsId() uint32 {
	return RTMP_CID_ProtocolControl
}
