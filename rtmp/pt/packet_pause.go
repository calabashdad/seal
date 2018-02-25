package pt

import (
	"fmt"
)

// PausePacket The client sends the pause command to tell the server to pause or start playing.
type PausePacket struct {

	// CommandName Name of the command, set to “pause”.
	CommandName string

	// TransactionID There is no transaction ID for this command. Set to 0.
	TransactionID float64

	// CommandObject Command information object does not exist. Set to null type.
	CommandObject Amf0Object

	// IsPause true or false, to indicate pausing or resuming play
	IsPause bool

	// TimeMs Number of milliseconds at which the the stream is paused or play resumed.
	// This is the current stream time at the Client when stream was paused. When the
	// playback is resumed, the server will only send messages with timestamps
	// greater than this value.
	TimeMs float64
}

// Decode .
func (pkt *PausePacket) Decode(data []uint8) (err error) {

	var offset uint32

	if pkt.CommandName, err = Amf0ReadString(data, &offset); err != nil {
		return
	}

	if RTMP_AMF0_COMMAND_PAUSE == pkt.CommandName {
		err = fmt.Errorf("decode pause packet command name is error.actully=%s", pkt.CommandName)
		return
	}

	if pkt.TransactionID, err = Amf0ReadNumber(data, &offset); err != nil {
		return
	}

	if err = amf0ReadNull(data, &offset); err != nil {
		return
	}

	if pkt.IsPause, err = amf0ReadBool(data, &offset); err != nil {
		return
	}

	if pkt.TimeMs, err = Amf0ReadNumber(data, &offset); err != nil {
		return
	}

	return
}

// Encode .
func (pkt *PausePacket) Encode() (data []uint8) {
	//no this method
	return
}

// GetMessageType .
func (pkt *PausePacket) GetMessageType() uint8 {
	//no this method
	return 0
}

// GetPreferCsID .
func (pkt *PausePacket) GetPreferCsID() uint32 {
	//no this method
	return 0
}
