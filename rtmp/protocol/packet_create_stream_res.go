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

func (pkt *CreateStreamResPacket) Encode() (data []uint8) {

	data = append(data, Amf0WriteString(pkt.Command_name)...)
	data = append(data, Amf0WriteNumber(pkt.Transaction_id)...)
	data = append(data, Amf0WriteNull()...)
	data = append(data, Amf0WriteNumber(pkt.Stream_id)...)

	return
}

func (pkt *CreateStreamResPacket) Decode(data []uint8) (err error) {

	var offset uint32

	err, pkt.Command_name = Amf0ReadString(data, &offset)
	if err != nil {
		return
	}

	if RTMP_AMF0_COMMAND_RESULT != pkt.Command_name {
		err = fmt.Errorf("decode create stream res packet, command name is not result.")
		return
	}

	err, pkt.Transaction_id = Amf0ReadNumber(data, &offset)
	if err != nil {
		return
	}

	err = Amf0ReadNull(data, &offset)
	if err != nil {
		return
	}

	err, pkt.Transaction_id = Amf0ReadNumber(data, &offset)
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
