package protocol

type CallResPacket struct {
}

func (pkt *CallResPacket) Decode(data []uint8) (err error) {
	return
}

func (pkt *CallResPacket) Encode() (data []uint8) {
	return
}

func (pkt *CallResPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (pkt *CallResPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverConnection
}
