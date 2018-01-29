package protocol

type SampleAccessPacket struct {

}

func (pkt* SampleAccessPacket)Decode([]uint8) (err error){
	return
}
func (pkt* SampleAccessPacket)Encode() (b []uint8){
	return
}
func (pkt* SampleAccessPacket)GetMessageType() uint8{
	return RTMP_MSG_AMF0DataMessage
}
func (pkt* SampleAccessPacket)GetPreferCsId() uint32{
	return RTMP_CID_OverStream
}