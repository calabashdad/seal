package pt

type CallPacket struct {
	/**
	 * Name of the remote procedure that is called.
	 */
	Command_name string
	/**
	 * If a response is expected we give a transaction Id. Else we pass a value of 0
	 */
	Transaction_id float64
	/**
	 * If there exists any command info this
	 * is set, else this is set to null type.
	 */
	Command_object  interface{}
	Cmd_object_type uint8
	/**
	 * Any optional arguments to be provided
	 */
	Arguments      interface{}
	Arguments_type uint8
}

func (pkt *CallPacket) Decode(data []uint8) (err error) {
	var offset uint32

	err, pkt.Command_name = Amf0ReadString(data, &offset)
	if err != nil {
		return
	}

	err, pkt.Transaction_id = Amf0ReadNumber(data, &offset)
	if err != nil {
		return
	}

	err, pkt.Command_object = Amf0ReadAny(data, &pkt.Cmd_object_type, &offset)
	if err != nil {
		return
	}

	if uint32(len(data))-offset > 0 {
		err, pkt.Arguments = Amf0ReadAny(data, &pkt.Arguments_type, &offset)
		if err != nil {
			return
		}
	}

	return
}

func (pkt *CallPacket) Encode() (data []uint8) {
	data = append(data, Amf0WriteString(pkt.Command_name)...)
	data = append(data, Amf0WriteNumber(pkt.Transaction_id)...)

	if nil != pkt.Command_object {
		data = append(data, Amf0WriteAny(pkt.Command_object.(Amf0Object))...)
	}

	if nil != pkt.Arguments {
		data = append(data, Amf0WriteAny(pkt.Arguments.(Amf0Object))...)
	}

	return
}

func (pkt *CallPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (pkt *CallPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverConnection
}
