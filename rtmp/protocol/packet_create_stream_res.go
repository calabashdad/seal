package protocol

import (
	"fmt"
)

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

	b = append(b, Amf0WriteString(pkt.Command_name)...)
	b = append(b, Amf0WriteNumber(pkt.Transaction_id)...)
	b = append(b, Amf0WriteNull()...)
	b = append(b, Amf0WriteNumber(pkt.Stream_id)...)

	return
}

func (pkt *CreateStreamResPacket) Decode(b []uint8) (err error) {

	var offset uint32

	err, pkt.Command_name = Amf0ReadString(b, &offset)
	if err != nil {
		return
	}

	if RTMP_AMF0_COMMAND_RESULT != pkt.Command_name {
		err = fmt.Errorf("decode create stream packet, command name is not result.")
		return
	}

	err, pkt.Transaction_id = Amf0ReadNumber(b, &offset)
	if err != nil {
		return
	}

	err = Amf0ReadNull(b, &offset)
	if err != nil {
		return
	}

	err, pkt.Transaction_id = Amf0ReadNumber(b, &offset)
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
