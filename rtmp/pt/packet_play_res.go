package pt

type PlayResPacket struct {
	/**
	 * Name of the command. If the play command is successful, the command
	 * name is set to onStatus.
	 */
	CommandName string
	/**
	 * Transaction ID set to 0.
	 */
	TransactionId float64
	/**
	 * Command information does not exist. Set to null type.
	 */
	CommandObject Amf0Object // null
	/**
	 * If the play command is successful, the client receives OnStatus message from
	 * server which is NetStream.Play.Start. If the specified stream is not found,
	 * NetStream.Play.StreamNotFound is received.
	 */
	Desc []Amf0Object
}

func (pkt *PlayResPacket) Decode(data []uint8) (err error) {
	//nothing
	return
}
func (pkt *PlayResPacket) Encode() (data []uint8) {
	data = append(data, Amf0WriteString(pkt.CommandName)...)
	data = append(data, Amf0WriteNumber(pkt.TransactionId)...)
	data = append(data, Amf0WriteObject(pkt.Desc)...)

	return
}
func (pkt *PlayResPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}
func (pkt *PlayResPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverStream
}
