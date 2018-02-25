package pt

// BandWidthPacket the special packet for the bandwidth test.
// actually, it's a OnStatusCallPacket, but
// 1. encode with data field, to send data to client.
// 2. decode ignore the data field, donot care.
type BandWidthPacket struct {
	// Name of command.
	CommandName string

	// Transaction ID set to 0.
	TransactionID float64

	// Args Command information does not exist. Set to null type.
	Args Amf0Object // null

	// Data Name-value pairs that describe the response from the server.
	// ‘code’,‘level’, ‘description’ are names of few among such information.
	Data []Amf0Object
}

// Decode .
func (pkt *BandWidthPacket) Decode(data []uint8) (err error) {
	var offset uint32

	if pkt.CommandName, err = Amf0ReadString(data, &offset); err != nil {
		return
	}

	if pkt.TransactionID, err = Amf0ReadNumber(data, &offset); err != nil {
		return
	}

	if err = amf0ReadNull(data, &offset); err != nil {
		return
	}

	if pkt.isStopPlay() || pkt.isStopPublish() || pkt.isFinish() {
		pkt.Data, err = amf0ReadObject(data, &offset)
		if err != nil {
			return
		}
	}

	return
}

func (pkt *BandWidthPacket) isStopPlay() bool {
	return SRS_BW_CHECK_STOP_PLAY == pkt.CommandName
}

func (pkt *BandWidthPacket) isStopPublish() bool {
	return SRS_BW_CHECK_START_PUBLISH == pkt.CommandName
}

func (pkt *BandWidthPacket) isFinish() bool {
	return SRS_BW_CHECK_FINISHED == pkt.CommandName
}

// Encode .
func (pkt *BandWidthPacket) Encode() (data []uint8) {
	data = append(data, amf0WriteString(pkt.CommandName)...)
	data = append(data, amf0WriteNumber(pkt.TransactionID)...)
	data = append(data, amf0WriteNull()...)
	data = append(data, amf0WriteObject(pkt.Data)...)

	return
}

// GetMessageType .
func (pkt *BandWidthPacket) GetMessageType() uint8 {
	return RtmpMsgAmf0CommandMessage
}

// GetPreferCsID .
func (pkt *BandWidthPacket) GetPreferCsID() uint32 {
	return RtmpCidOverStream
}
