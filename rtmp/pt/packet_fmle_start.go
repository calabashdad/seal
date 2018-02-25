package pt

import (
	"fmt"
	"strings"
)

// FmleStartPacket FMLE start publish: ReleaseStream/PublishStream
type FmleStartPacket struct {
	// CommandName Name of the command
	CommandName string

	// TransactionID the transaction ID to get the response.
	TransactionID float64

	// CommandObject If there exists any command info this is set, else this is set to null type.
	CommandObject Amf0Object // null

	// StreamName the stream name to start publish or release.
	StreamName string

	// TokenStrToken value, for authentication. it's optional.
	TokenStr string
}

// Decode .
func (pkt *FmleStartPacket) Decode(data []uint8) (err error) {
	var offset uint32

	if pkt.CommandName, err = Amf0ReadString(data, &offset); err != nil {
		return
	}

	if RTMP_AMF0_COMMAND_RELEASE_STREAM != pkt.CommandName &&
		RTMP_AMF0_COMMAND_FC_PUBLISH != pkt.CommandName &&
		RTMP_AMF0_COMMAND_UNPUBLISH != pkt.CommandName {
		err = fmt.Errorf("decode fmle start packet error, command name is error.actully=%s", pkt.CommandName)
		return
	}

	if pkt.TransactionID, err = Amf0ReadNumber(data, &offset); err != nil {
		return
	}

	if err = amf0ReadNull(data, &offset); err != nil {
		return
	}

	var streamNameLocal string
	if streamNameLocal, err = Amf0ReadString(data, &offset); err != nil {
		return
	}

	i := strings.Index(streamNameLocal, TokenStr)
	if i < 0 {
		pkt.StreamName = streamNameLocal
	} else {
		pkt.StreamName = streamNameLocal[0:i]
		pkt.TokenStr = streamNameLocal[i+len(TokenStr):]
	}

	return
}

// Encode .
func (pkt *FmleStartPacket) Encode() (data []uint8) {
	data = append(data, amf0WriteString(pkt.CommandName)...)
	data = append(data, amf0WriteNumber(pkt.TransactionID)...)
	data = append(data, amf0WriteNull()...)
	data = append(data, amf0WriteString(pkt.StreamName)...)

	return
}

// GetMessageType .
func (pkt *FmleStartPacket) GetMessageType() uint8 {
	return RtmpMsgAmf0CommandMessage
}

// GetPreferCsID .
func (pkt *FmleStartPacket) GetPreferCsID() uint32 {
	return RtmpCidOverConnection
}
