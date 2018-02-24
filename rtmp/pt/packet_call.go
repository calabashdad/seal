package pt

type CallPacket struct {
	/**
	 * Name of the remote procedure that is called.
	 */
	CommandName string
	/**
	 * If a response is expected we give a transaction Id. Else we pass a value of 0
	 */
	TransactionId float64
	/**
	 * If there exists any command info this
	 * is set, else this is set to null type.
	 */
	CommandObject  interface{}
	Cmd_objectType uint8
	/**
	 * Any optional arguments to be provided
	 */
	Arguments     interface{}
	ArgumentsType uint8
}

func (pkt *CallPacket) Decode(data []uint8) (err error) {
	var offset uint32

	pkt.CommandName, err = Amf0ReadString(data, &offset)
	if err != nil {
		return
	}

	pkt.TransactionId, err = Amf0ReadNumber(data, &offset)
	if err != nil {
		return
	}

	pkt.CommandObject, err = amf0ReadAny(data, &pkt.Cmd_objectType, &offset)
	if err != nil {
		return
	}

	if uint32(len(data))-offset > 0 {
		pkt.Arguments, err = amf0ReadAny(data, &pkt.ArgumentsType, &offset)
		if err != nil {
			return
		}
	}

	return
}

func (pkt *CallPacket) Encode() (data []uint8) {
	data = append(data, amf0WriteString(pkt.CommandName)...)
	data = append(data, amf0WriteNumber(pkt.TransactionId)...)

	if nil != pkt.CommandObject {
		data = append(data, amf0WriteAny(pkt.CommandObject.(Amf0Object))...)
	}

	if nil != pkt.Arguments {
		data = append(data, amf0WriteAny(pkt.Arguments.(Amf0Object))...)
	}

	return
}

func (pkt *CallPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (pkt *CallPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverConnection
}
