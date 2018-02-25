package pt

import "fmt"

// FmleStartResPacket response for FmleStartPacket
type FmleStartResPacket struct {

	// CommandName Name of the command
	CommandName string

	// TransactionID the transaction ID to get the response.
	TransactionID float64

	// CommandObject If there exists any command info this is set, else this is set to null type.
	CommandObject Amf0Object // null

	// Args the optional args, set to undefined.
	Args Amf0Object // undefined
}

// Decode .
func (pkt *FmleStartResPacket) Decode(data []uint8) (err error) {
	var offset uint32

	if pkt.CommandName, err = Amf0ReadString(data, &offset); err != nil {
		return
	}

	if RTMP_AMF0_COMMAND_RESULT != pkt.CommandName {
		err = fmt.Errorf("decode fmle start res packet, command name is not result")
		return
	}

	if pkt.TransactionID, err = Amf0ReadNumber(data, &offset); err != nil {
		return
	}

	if err = amf0ReadNull(data, &offset); err != nil {
		return
	}

	if err = amf0ReadUndefined(data, &offset); err != nil {
		return
	}

	return
}

// Encode .
func (pkt *FmleStartResPacket) Encode() (data []uint8) {

	data = append(data, amf0WriteString(pkt.CommandName)...)
	data = append(data, amf0WriteNumber(pkt.TransactionID)...)
	data = append(data, amf0WriteNull()...)
	data = append(data, amf0WriteUndefined()...)

	return
}

// GetMessageType .
func (pkt *FmleStartResPacket) GetMessageType() uint8 {
	return RtmpMsgAmf0CommandMessage
}

// GetPreferCsID .
func (pkt *FmleStartResPacket) GetPreferCsID() uint32 {
	return RtmpCidOverConnection
}
