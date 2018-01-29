package protocol

/**
* FMLE start publish: ReleaseStream/PublishStream
 */
type FmleStartPacket struct {
}

func (pkt *FmleStartPacket) Decode([]uint8) (err error) {
	return
}
func (pkt *FmleStartPacket) Encode() (b []uint8) {
	return
}
func (pkt *FmleStartPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}
func (pkt *FmleStartPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverConnection
}
