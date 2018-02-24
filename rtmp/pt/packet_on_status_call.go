package pt

type OnStatusCallPacket struct {
	/**
	 * Name of command. Set to "onStatus"
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
	/**
	 * Name-value pairs that describe the response from the server.
	 * ‘code’,‘level’, ‘description’ are names of few among such information.
	 */
	Data []Amf0Object
}

func (pkt *OnStatusCallPacket) Decode(data []uint8) (err error) {
	//nothing
	return
}
func (pkt *OnStatusCallPacket) Encode() (data []uint8) {
	data = append(data, amf0WriteString(pkt.CommandName)...)
	data = append(data, amf0WriteNumber(pkt.TransactionId)...)
	data = append(data, amf0WriteNull()...)
	data = append(data, amf0WriteObject(pkt.Data)...)

	return
}
func (pkt *OnStatusCallPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}
func (pkt *OnStatusCallPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverStream
}
