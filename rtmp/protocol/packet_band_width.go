package protocol

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
	Command_name string
	/**
	 * Transaction ID set to 0.
	 */
	Transaction_id float64
	/**
	 * Command information does not exist. Set to null type.
	 */
	Args Amf0Object // null
	/**
	 * Name-value pairs that describe the response from the server.
	 * ‘code’,‘level’, ‘description’ are names of few among such information.
	 */
	Data Amf0Object
}

func (pkt *BandWidthPacket) Decode(data []uint8) (err error) {
	return
}

func (pkt *BandWidthPacket) Encode() (data []uint8) {
	return
}

func (pkt *BandWidthPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (pkt *BandWidthPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverStream
}
