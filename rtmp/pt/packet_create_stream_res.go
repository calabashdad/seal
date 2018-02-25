package pt

import (
	"fmt"
)

// CreateStreamResPacket response for CreateStreamPacket
type CreateStreamResPacket struct {

	// CommandName  _result or _error; indicates whether the response is result or error.
	CommandName string

	// TransactionID ID of the command that response belongs to.
	TransactionID float64

	// CommandObject If there exists any command info this is set, else this is set to null type.
	CommandObject Amf0Object // null

	// StreamID The return value is either a stream ID or an error information object.
	StreamID float64
}

// Encode .
func (pkt *CreateStreamResPacket) Encode() (data []uint8) {

	data = append(data, amf0WriteString(pkt.CommandName)...)
	data = append(data, amf0WriteNumber(pkt.TransactionID)...)
	data = append(data, amf0WriteNull()...)
	data = append(data, amf0WriteNumber(pkt.StreamID)...)

	return
}

// Decode .
func (pkt *CreateStreamResPacket) Decode(data []uint8) (err error) {

	var offset uint32

	if pkt.CommandName, err = Amf0ReadString(data, &offset); err != nil {
		return
	}

	if RTMP_AMF0_COMMAND_RESULT != pkt.CommandName {
		err = fmt.Errorf("decode create stream res packet, command name is not result")
		return
	}

	if pkt.TransactionID, err = Amf0ReadNumber(data, &offset); err != nil {
		return
	}

	if err = amf0ReadNull(data, &offset); err != nil {
		return
	}

	if pkt.TransactionID, err = Amf0ReadNumber(data, &offset); err != nil {
		return
	}

	return
}

// GetMessageType .
func (pkt *CreateStreamResPacket) GetMessageType() uint8 {
	return RtmpMsgAmf0CommandMessage
}

// GetPreferCsID .
func (pkt *CreateStreamResPacket) GetPreferCsID() uint32 {
	return RtmpCidOverConnection
}
