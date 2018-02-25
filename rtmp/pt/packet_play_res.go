package pt

// PlayResPacket response for PlayPacket
type PlayResPacket struct {
	// CommandName Name of the command. If the play command is successful, the command
	// name is set to onStatus.
	CommandName string

	// Transaction ID set to 0.
	TransactionID float64

	// CommandObject Command information does not exist. Set to null type.
	CommandObject Amf0Object

	// Desc If the play command is successful, the client receives OnStatus message from
	// server which is NetStream.Play.Start. If the specified stream is not found,
	// NetStream.Play.StreamNotFound is received.
	Desc []Amf0Object
}

// Decode .
func (pkt *PlayResPacket) Decode(data []uint8) (err error) {
	//nothing
	return
}

// Encode .
func (pkt *PlayResPacket) Encode() (data []uint8) {
	data = append(data, amf0WriteString(pkt.CommandName)...)
	data = append(data, amf0WriteNumber(pkt.TransactionID)...)
	data = append(data, amf0WriteObject(pkt.Desc)...)

	return
}

// GetMessageType .
func (pkt *PlayResPacket) GetMessageType() uint8 {
	return RtmpMsgAmf0CommandMessage
}

// GetPreferCsID .
func (pkt *PlayResPacket) GetPreferCsID() uint32 {
	return RtmpCidOverStream
}
