package protocol

type OnStatusDataPacket struct {

}

func (pkt* OnStatusDataPacket)Decode([]uint8) (err error){
	return
}
func (pkt* OnStatusDataPacket)Encode() (b []uint8){
	return
}
func (pkt* OnStatusDataPacket)GetMessageType() uint8{
	return RTMP_MSG_AMF0DataMessage
}
func (pkt* OnStatusDataPacket)GetPreferCsId() uint32{
	return RTMP_CID_OverStream
}