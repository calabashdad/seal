package protocol

type CreateStreamResPacket struct {
	/**
	 * _result or _error; indicates whether the response is result or error.
	 */
	Command_name string

	/**
	 * ID of the command that response belongs to.
	 */
	Transaction_id float64
	/**
	 * If there exists any command info this is set, else this is set to null type.
	 */
	Command_object Amf0Object // null
	/**
	 * The return value is either a stream ID or an error information object.
	 */
	Stream_id float64
}

func (pkt *CreateStreamResPacket) Encode() (b []uint8) {
	return
}

func (pkt *CreateStreamResPacket) Decode(b []uint8) (err error) {
	return
}

func (pkt *CreateStreamResPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (pkt *CreateStreamResPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverConnection
}
