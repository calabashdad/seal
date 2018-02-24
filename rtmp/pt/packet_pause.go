package pt

import (
	"fmt"
)

type PausePacket struct {
	/**
	 * Name of the command, set to “pause”.
	 */
	CommandName string
	/**
	 * There is no transaction ID for this command. Set to 0.
	 */
	TransactionId float64
	/**
	 * Command information object does not exist. Set to null type.
	 */
	CommandObject Amf0Object // null
	/**
	 * true or false, to indicate pausing or resuming play
	 */
	IsPause bool
	/**
	 * Number of milliseconds at which the the stream is paused or play resumed.
	 * This is the current stream time at the Client when stream was paused. When the
	 * playback is resumed, the server will only send messages with timestamps
	 * greater than this value.
	 */
	TimeMs float64
}

func (pkt *PausePacket) Decode(data []uint8) (err error) {

	var offset uint32

	pkt.CommandName, err = Amf0ReadString(data, &offset)
	if err != nil {
		return
	}

	if RTMP_AMF0_COMMAND_PAUSE == pkt.CommandName {
		err = fmt.Errorf("decode pause packet command name is error.actully=%s", pkt.CommandName)
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

	pkt.IsPause, err = amf0ReadBool(data, &offset)
	if err != nil {
		return
	}

	pkt.TimeMs, err = Amf0ReadNumber(data, &offset)
	if err != nil {
		return
	}

	return
}

func (pkt *PausePacket) Encode() (data []uint8) {
	//no this method

	return
}

func (pkt *PausePacket) GetMessageType() uint8 {
	//no this method
	return 0
}
func (pkt *PausePacket) GetPreferCsId() uint32 {
	//no this method

	return 0
}
