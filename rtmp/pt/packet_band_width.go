package pt

/**
* the special packet for the bandwidth test.
* actually, it's a OnStatusCallPacket, but
* 1. encode with data field, to send data to client.
* 2. decode ignore the data field, donot care.
 */
type BandWidthPacket struct {
	/**
	 * Name of command.
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

func (pkt *BandWidthPacket) Decode(data []uint8) (err error) {
	var offset uint32

	pkt.CommandName, err = Amf0ReadString(data, &offset)
	if err != nil {
		return
	}

	pkt.TransactionId, err = Amf0ReadNumber(data, &offset)
	if err != nil {
		return
	}

	err = Amf0ReadNull(data, &offset)
	if err != nil {
		return
	}

	if pkt.isStopPlay() || pkt.isStopPublish() || pkt.isFinish() {
		pkt.Data, err = Amf0ReadObject(data, &offset)
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
func (pkt *BandWidthPacket) Encode() (data []uint8) {
	data = append(data, Amf0WriteString(pkt.CommandName)...)
	data = append(data, Amf0WriteNumber(pkt.TransactionId)...)
	data = append(data, Amf0WriteNull()...)
	data = append(data, Amf0WriteObject(pkt.Data)...)

	return
}

func (pkt *BandWidthPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (pkt *BandWidthPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverStream
}
