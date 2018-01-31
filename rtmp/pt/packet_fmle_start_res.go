package pt

import "fmt"

type FmleStartResPacket struct {
	/**
	 * Name of the command
	 */
	CommandName string
	/**
	 * the transaction ID to get the response.
	 */
	TransactionId float64
	/**
	 * If there exists any command info this is set, else this is set to null type.
	 */
	CommandObject Amf0Object // null
	/**
	 * the optional args, set to undefined.
	 */
	Args Amf0Object // undefined
}

func (pkt *FmleStartResPacket) Decode(data []uint8) (err error) {
	var offset uint32

	err, pkt.CommandName = Amf0ReadString(data, &offset)
	if err != nil {
		return
	}

	if RTMP_AMF0_COMMAND_RESULT != pkt.CommandName {
		err = fmt.Errorf("decode fmle start res packet, command name is not result.")
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

	err = Amf0ReadUndefined(data, &offset)
	if err != nil {
		return
	}

	return
}

func (pkt *FmleStartResPacket) Encode() (data []uint8) {

	data = append(data, Amf0WriteString(pkt.CommandName)...)
	data = append(data, Amf0WriteNumber(pkt.TransactionId)...)
	data = append(data, Amf0WriteNull()...)
	data = append(data, Amf0WriteUndefined()...)

	return
}

func (pkt *FmleStartResPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (pkt *FmleStartResPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverConnection
}
