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

	pkt.CommandName, err = Amf0ReadString(data, &offset)
	if err != nil {
		return
	}

	if RTMP_AMF0_COMMAND_RESULT != pkt.CommandName {
		err = fmt.Errorf("decode fmle start res packet, command name is not result.")
		return
	}

	pkt.TransactionId, err = Amf0ReadNumber(data, &offset)
	if err != nil {
		return
	}

	err = amf0ReadNull(data, &offset)
	if err != nil {
		return
	}

	err = amf0ReadUndefined(data, &offset)
	if err != nil {
		return
	}

	return
}

func (pkt *FmleStartResPacket) Encode() (data []uint8) {

	data = append(data, amf0WriteString(pkt.CommandName)...)
	data = append(data, amf0WriteNumber(pkt.TransactionId)...)
	data = append(data, amf0WriteNull()...)
	data = append(data, amf0WriteUndefined()...)

	return
}

func (pkt *FmleStartResPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (pkt *FmleStartResPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverConnection
}
