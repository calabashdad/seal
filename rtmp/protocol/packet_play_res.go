package protocol

type PlayResPacket struct {

}

func (pkt* PlayResPacket)Decode([]uint8) (err error){
	return
}
func (pkt* PlayResPacket)Encode() (b []uint8){
	return
}
func (pkt* PlayResPacket)GetMessageType() uint8{
	return RTMP_MSG_AMF0CommandMessage
}
func (pkt* PlayResPacket)GetPreferCsId() uint32{
	return RTMP_CID_OverStream
}