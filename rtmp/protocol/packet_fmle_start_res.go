package protocol

type FmleStartResPacket struct {
	/**
	 * Name of the command
	 */
	Command_name string
	/**
	 * the transaction ID to get the response.
	 */
	Transaction_id float64
	/**
	 * If there exists any command info this is set, else this is set to null type.
	 */
	Command_object Amf0Object // null
	/**
	 * the optional args, set to undefined.
	 */
	Args Amf0Object // undefined
}

func (pkt *FmleStartResPacket) Decode([]uint8) (err error) {

	return
}

func (pkt *FmleStartResPacket) Encode() (b []uint8) {

	return
}

func (pkt *FmleStartResPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (pkt *FmleStartResPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverConnection
}
