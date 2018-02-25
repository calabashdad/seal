package pt

// CloseStreamPacket client close stream packet.
type CloseStreamPacket struct {
	// CommandName Name of the command, set to “closeStream”.
	CommandName string

	// Transaction ID set to 0.
	TransactionID float64

	// CommandObject Command information object does not exist. Set to null type.
	CommandObject Amf0Object // null
}

// Decode .
func (pkt *CloseStreamPacket) Decode(data []uint8) (err error) {
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

	return
}

// Encode .
func (pkt *CloseStreamPacket) Encode() (data []uint8) {
	//nothing
	return
}

// GetMessageType .
func (pkt *CloseStreamPacket) GetMessageType() uint8 {
	//no method for this pakcet
	return 0
}

// GetPreferCsID .
func (pkt *CloseStreamPacket) GetPreferCsID() uint32 {
	//no method for this pakcet
	return 0
}
