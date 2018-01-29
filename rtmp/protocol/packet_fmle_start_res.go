package protocol

import "fmt"

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

func (pkt *FmleStartResPacket) Decode(b []uint8) (err error) {
	var offset uint32

	err, pkt.Command_name = Amf0ReadString(b, &offset)
	if err != nil {
		return
	}

	if RTMP_AMF0_COMMAND_RESULT != pkt.Command_name {
		err = fmt.Errorf("decode fmle start res packet, command name is not result.")
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

	err = Amf0ReadUndefined(b, &offset)
	if err != nil {
		return
	}

	return
}

func (pkt *FmleStartResPacket) Encode() (b []uint8) {

	b = append(b, Amf0WriteString(pkt.Command_name)...)
	b = append(b, Amf0WriteNumber(pkt.Transaction_id)...)
	b = append(b, Amf0WriteNull()...)
	b = append(b, Amf0WriteUndefined()...)

	return
}

func (pkt *FmleStartResPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (pkt *FmleStartResPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverConnection
}
