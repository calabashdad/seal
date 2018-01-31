package pt

type CloseStreamPacket struct {
	/**
	 * Name of the command, set to “closeStream”.
	 */
	CommandName string
	/**
	 * Transaction ID set to 0.
	 */
	TransactionId float64
	/**
	 * Command information object does not exist. Set to null type.
	 */
	CommandObject Amf0Object // null
}

func (pkt *CloseStreamPacket) Decode(data []uint8) (err error) {
	var offset uint32

	err, pkt.CommandName = Amf0ReadString(data, &offset)
	if err != nil {
		return
	}

	err, pkt.TransactionId = Amf0ReadNumber(data, &offset)
	if err != nil {
		return
	}

	err = Amf0ReadNull(data, &offset)
	if err != nil {
		return
	}

	return
}
func (pkt *CloseStreamPacket) Encode() (data []uint8) {
	//nothing
	return
}
func (pkt *CloseStreamPacket) GetMessageType() uint8 {
	//no method for this pakcet
	return 0
}
func (pkt *CloseStreamPacket) GetPreferCsId() uint32 {
	//no method for this pakcet
	return 0
}
