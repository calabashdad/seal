package protocol

type SetWindowAckSizePacket struct {

}

func (pkt* SetWindowAckSizePacket)Decode([]uint8) (err error){
	return
}
func (pkt* SetWindowAckSizePacket)Encode() (b []uint8){
	return
}
func (pkt* SetWindowAckSizePacket)GetMessageType() uint8{
	return RTMP_MSG_WindowAcknowledgementSize
}
func (pkt* SetWindowAckSizePacket)GetPreferCsId() uint32{
	return RTMP_CID_ProtocolControl
}