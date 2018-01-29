package protocol

type SetChunkSizePacket struct {
}

func (pkt *SetChunkSizePacket) Decode(data []uint8) (err error) {
	return
}
func (pkt *SetChunkSizePacket) Encode() (data []uint8) {
	return
}
func (pkt *SetChunkSizePacket) GetMessageType() uint8 {
	return RTMP_MSG_SetChunkSize
}
func (pkt *SetChunkSizePacket) GetPreferCsId() uint32 {
	return RTMP_CID_ProtocolControl
}
