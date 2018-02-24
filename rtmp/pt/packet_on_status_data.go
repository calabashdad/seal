package pt

type OnStatusDataPacket struct {
	/**
	 * Name of command. Set to "onStatus"
	 */
	CommandName string
	/**
	 * Name-value pairs that describe the response from the server.
	 * ‘code’, are names of few among such information.
	 */
	Data []Amf0Object
}

func (pkt *OnStatusDataPacket) Decode(data []uint8) (err error) {
	//nothing
	return
}
func (pkt *OnStatusDataPacket) Encode() (data []uint8) {
	data = append(data, amf0WriteString(pkt.CommandName)...)
	data = append(data, amf0WriteObject(pkt.Data)...)

	return
}
func (pkt *OnStatusDataPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0DataMessage
}
func (pkt *OnStatusDataPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverStream
}
