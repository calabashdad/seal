package protocol

type SetPeerBandWidthPacket struct {
}

func (pkt *SetPeerBandWidthPacket) Decode(data []uint8) (err error) {
	return
}
func (pkt *SetPeerBandWidthPacket) Encode() (data []uint8) {
	return
}
func (pkt *SetPeerBandWidthPacket) GetMessageType() uint8 {
	return RTMP_MSG_SetPeerBandwidth
}
func (pkt *SetPeerBandWidthPacket) GetPreferCsId() uint32 {
	return RTMP_CID_ProtocolControl
}
