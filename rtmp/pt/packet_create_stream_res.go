package pt

import (
	"fmt"
)

type CreateStreamResPacket struct {
	/**
	 * _result or _error; indicates whether the response is result or error.
	 */
	CommandName string

	/**
	 * ID of the command that response belongs to.
	 */
	TransactionId float64
	/**
	 * If there exists any command info this is set, else this is set to null type.
	 */
	CommandObject Amf0Object // null
	/**
	 * The return value is either a stream ID or an error information object.
	 */
	StreamId float64
}

func (pkt *CreateStreamResPacket) Encode() (data []uint8) {

	data = append(data, amf0WriteString(pkt.CommandName)...)
	data = append(data, amf0WriteNumber(pkt.TransactionId)...)
	data = append(data, amf0WriteNull()...)
	data = append(data, amf0WriteNumber(pkt.StreamId)...)

	return
}

func (pkt *CreateStreamResPacket) Decode(data []uint8) (err error) {

	var offset uint32

	pkt.CommandName, err = Amf0ReadString(data, &offset)
	if err != nil {
		return
	}

	if RTMP_AMF0_COMMAND_RESULT != pkt.CommandName {
		err = fmt.Errorf("decode create stream res packet, command name is not result.")
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

	pkt.TransactionId, err = Amf0ReadNumber(data, &offset)
	if err != nil {
		return
	}

	return
}

func (pkt *CreateStreamResPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (pkt *CreateStreamResPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverConnection
}
