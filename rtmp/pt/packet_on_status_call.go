package pt

// OnStatusCallPacket onStatus command, AMF0 Call
// user must set the stream id
type OnStatusCallPacket struct {

	// CommandName Name of command. Set to "onStatus"
	CommandName string

	// TransactionID set to 0
	TransactionID float64

	// Args Command information does not exist. Set to null type.
	Args Amf0Object

	// Data Name-value pairs that describe the response from the server.
	// ‘code’,‘level’, ‘description’ are names of few among such information.
	Data []Amf0Object
}

// Decode .
func (pkt *OnStatusCallPacket) Decode(data []uint8) (err error) {
	//nothing
	return
}

// Encode .
func (pkt *OnStatusCallPacket) Encode() (data []uint8) {
	data = append(data, amf0WriteString(pkt.CommandName)...)
	data = append(data, amf0WriteNumber(pkt.TransactionID)...)
	data = append(data, amf0WriteNull()...)
	data = append(data, amf0WriteObject(pkt.Data)...)

	return
}

// GetMessageType .
func (pkt *OnStatusCallPacket) GetMessageType() uint8 {
	return RtmpMsgAmf0CommandMessage
}

// GetPreferCsID .
func (pkt *OnStatusCallPacket) GetPreferCsID() uint32 {
	return RtmpCidOverStream
}

// AddObj add object to data
func (pkt *OnStatusCallPacket) AddObj(obj *Amf0Object) {
	pkt.Data = append(pkt.Data, *obj)
}
