package pt

import (
	"fmt"
)

// CreateStreamPacket  The client sends this command to the server to create a logical
// channel for message communication The publishing of audio, video, and
// metadata is carried out over stream channel created using the
// createStream command.
type CreateStreamPacket struct {
	// CommandName Name of the command. Set to “createStream”.
	CommandName string

	// TransactionID Transaction ID of the command.
	TransactionID float64

	// CommandObject If there exists any command info this is set, else this is set to null type.
	CommandObject Amf0Object // null
}

// Decode .
func (pkt *CreateStreamPacket) Decode(data []uint8) (err error) {
	var offset uint32

	if pkt.CommandName, err = Amf0ReadString(data, &offset); err != nil {
		return
	}

	if RTMP_AMF0_COMMAND_CREATE_STREAM != pkt.CommandName {
		err = fmt.Errorf("decode create stream packet, command name is wrong. actully=%s", pkt.CommandName)
		return
	}

	if pkt.TransactionID, err = Amf0ReadNumber(data, &offset); err != nil {
		return
	}

	if err = amf0ReadNull(data, &offset); err != nil {
		return
	}

	return
}

// Encode .
func (pkt *CreateStreamPacket) Encode() (data []uint8) {
	data = append(data, amf0WriteString(pkt.CommandName)...)
	data = append(data, amf0WriteNumber(pkt.TransactionID)...)
	data = append(data, amf0WriteNull()...)

	return
}

// GetMessageType .
func (pkt *CreateStreamPacket) GetMessageType() uint8 {
	return RtmpMsgAmf0CommandMessage
}

// GetPreferCsID .
func (pkt *CreateStreamPacket) GetPreferCsID() uint32 {
	return RtmpCidOverConnection
}
