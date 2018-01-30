package pt

type OnBwDonePacket struct {
	/**
	 * Name of command. Set to "onBWDone"
	 */
	CommandName string
	/**
	 * Transaction ID set to 0.
	 */
	TransactionId float64
	/**
	 * Command information does not exist. Set to null type.
	 */
	Args Amf0Object // null
}

func (pkt *OnBwDonePacket) Decode(data []uint8) (err error) {
	return
}
func (pkt *OnBwDonePacket) Encode() (data []uint8) {
	data = append(data, Amf0WriteString(pkt.CommandName)...)
	data = append(data, Amf0WriteNumber(pkt.TransactionId)...)
	data = append(data, Amf0WriteNull()...)

	return
}
func (pkt *OnBwDonePacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}
func (pkt *OnBwDonePacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverConnection
}
