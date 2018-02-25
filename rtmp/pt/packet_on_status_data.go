package pt

// OnStatusDataPacket onStatus data, AMF0 Data
// user must set the stream id
type OnStatusDataPacket struct {

	// CommandName Name of command. Set to "onStatus"
	CommandName string

	// Data Name-value pairs that describe the response from the server.
	// ‘code’, are names of few among such information.
	Data []Amf0Object
}

// Decode .
func (pkt *OnStatusDataPacket) Decode(data []uint8) (err error) {
	//nothing
	return
}

// Encode .
func (pkt *OnStatusDataPacket) Encode() (data []uint8) {
	data = append(data, amf0WriteString(pkt.CommandName)...)
	data = append(data, amf0WriteObject(pkt.Data)...)

	return
}

// GetMessageType .
func (pkt *OnStatusDataPacket) GetMessageType() uint8 {
	return RtmpMsgAmf0DataMessage
}

// GetPreferCsID .
func (pkt *OnStatusDataPacket) GetPreferCsID() uint32 {
	return RtmpCidOverStream
}

// AddObj add object to Data
func (pkt *OnStatusDataPacket) AddObj(obj *Amf0Object) {
	pkt.Data = append(pkt.Data, *obj)
}
