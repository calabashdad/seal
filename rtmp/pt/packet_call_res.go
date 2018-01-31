package pt

type CallResPacket struct {
	/**
	 * Name of the command.
	 */
	CommandName string

	/**
	 * @brief 请求的命令名
	 */
	ReqCommandName string

	/**
	 * ID of the command, to which the response belongs to
	 */
	TransactionId float64
	/**
	 * If there exists any command info this is set, else this is set to null type.
	 */
	CommandObject       interface{}
	CommandObjectMarker uint8
	/**
	 * Response from the method that was called.
	 */
	Response       interface{}
	ResponseMarker uint8
}

func (pkt *CallResPacket) Decode(data []uint8) (err error) {
	var offset uint32

	err, pkt.CommandName = Amf0ReadString(data, &offset)
	if err != nil {
		return
	}

	err, pkt.TransactionId = Amf0ReadNumber(data, &offset)
	if err != nil {
		return
	}

	err, pkt.CommandObject = Amf0ReadAny(data, &pkt.CommandObjectMarker, &offset)
	if err != nil {
		return
	}

	maxOffset := uint32(len(data)) - 1
	if maxOffset-offset > 0 {
		err, pkt.Response = Amf0ReadAny(data, &pkt.ResponseMarker, &offset)
		if err != nil {
			return
		}

	}

	return
}

func (pkt *CallResPacket) Encode() (data []uint8) {
	data = append(data, Amf0WriteString(pkt.CommandName)...)
	data = append(data, Amf0WriteNumber(pkt.TransactionId)...)
	if nil != pkt.CommandObject {
		data = append(data, Amf0WriteAny(pkt.CommandObject.(Amf0Object))...)
	}

	if nil != pkt.Response {
		data = append(data, Amf0WriteAny(pkt.Response.(Amf0Object))...)
	}

	return
}

func (pkt *CallResPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (pkt *CallResPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverConnection
}
