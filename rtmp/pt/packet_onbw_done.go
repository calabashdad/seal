package pt

// OnBwDonePacket when bandwidth test done, notice client
type OnBwDonePacket struct {

	// CommandName Name of command. Set to "onBWDone"
	CommandName string

	// TransactionID Transaction ID set to 0.
	TransactionID float64

	// Args  Command information does not exist. Set to null type.
	Args Amf0Object
}

// Decode .
func (pkt *OnBwDonePacket) Decode(data []uint8) (err error) {
	//nothing
	return
}

// Encode .
func (pkt *OnBwDonePacket) Encode() (data []uint8) {
	data = append(data, amf0WriteString(pkt.CommandName)...)
	data = append(data, amf0WriteNumber(pkt.TransactionID)...)
	data = append(data, amf0WriteNull()...)

	return
}

// GetMessageType .
func (pkt *OnBwDonePacket) GetMessageType() uint8 {
	return RtmpMsgAmf0CommandMessage
}

// GetPreferCsID .
func (pkt *OnBwDonePacket) GetPreferCsID() uint32 {
	return RtmpCidOverConnection
}
