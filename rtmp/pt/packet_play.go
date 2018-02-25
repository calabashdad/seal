package pt

import (
	"fmt"
	"strings"
)

// PlayPacket The client sends this command to the server to play a stream.
type PlayPacket struct {

	// CommandName Name of the command. Set to “play”.
	CommandName string

	// Transaction ID set to 0.
	TransactionID float64

	// CommandObject Command information does not exist. Set to null type.
	CommandObject Amf0Object

	// StreamName Name of the stream to play.
	// To play video (FLV) files, specify the name of the stream without a file
	//       extension (for example, "sample").
	// To play back MP3 or ID3 tags, you must precede the stream name with mp3:
	//       (for example, "mp3:sample".)
	// To play H.264/AAC files, you must precede the stream name with mp4: and specify the
	//       file extension. For example, to play the file sample.m4v, specify
	//       "mp4:sample.m4v"
	StreamName string

	// TokenStr Token value, for authentication. it's optional.
	TokenStr string

	// Start An optional parameter that specifies the start time in seconds.
	// The default value is -2, which means the subscriber first tries to play the live
	//       stream specified in the Stream Name field. If a live stream of that name is
	//       not found, it plays the recorded stream specified in the Stream Name field.
	// If you pass -1 in the Start field, only the live stream specified in the Stream
	//       Name field is played.
	// If you pass 0 or a positive number in the Start field, a recorded stream specified
	//       in the Stream Name field is played beginning from the time specified in the
	//       Start field.
	// If no recorded stream is found, the next item in the playlist is played.
	Start float64

	// Duration An optional parameter that specifies the duration of playback in seconds.
	// The default value is -1. The -1 value means a live stream is played until it is no
	//       longer available or a recorded stream is played until it ends.
	// If u pass 0, it plays the single frame since the time specified in the Start field
	//       from the beginning of a recorded stream. It is assumed that the value specified
	//       in the Start field is equal to or greater than 0.
	// If you pass a positive number, it plays a live stream for the time period specified
	//       in the Duration field. After that it becomes available or plays a recorded
	//       stream for the time specified in the Duration field. (If a stream ends before the
	//       time specified in the Duration field, playback ends when the stream ends.)
	// If you pass a negative number other than -1 in the Duration field, it interprets the
	//       value as if it were -1.
	Duration float64

	// Reset An optional Boolean value or number that specifies whether to flush any
	// previous playlist.
	Reset bool
}

// Decode .
func (pkt *PlayPacket) Decode(data []uint8) (err error) {
	var maxOffset uint32
	maxOffset = uint32(len(data)) - 1

	var offset uint32

	if pkt.CommandName, err = Amf0ReadString(data, &offset); err != nil {
		return
	}

	if RTMP_AMF0_COMMAND_PLAY != pkt.CommandName {
		err = fmt.Errorf("decode play packet, command name is not play.actully=%s", pkt.CommandName)
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

	if maxOffset-offset > (1 + 8) { // number need at least 1(marker) + 8(number)
		if pkt.Start, err = Amf0ReadNumber(data, &offset); err != nil {
			return
		}
	}

	if maxOffset-offset > (1 + 8) {
		if pkt.Duration, err = Amf0ReadNumber(data, &offset); err != nil {
			return
		}
	}

	if offset >= uint32(len(data)) {
		return
	}

	if maxOffset-offset >= 2 { //because the bool type need 2 bytes at least

		var v interface{}
		var marker uint8
		if v, err = amf0ReadAny(data, &marker, &offset); err != nil {
			return
		}

		if RTMP_AMF0_Boolean == marker {
			pkt.Reset = v.(bool)
		} else if RTMP_AMF0_Number == marker {
			pkt.Reset = (v.(float64) != 0)
		}
	}

	return
}

// Encode .
func (pkt *PlayPacket) Encode() (data []uint8) {

	data = append(data, amf0WriteString(pkt.CommandName)...)
	data = append(data, amf0WriteNumber(pkt.TransactionID)...)
	data = append(data, amf0WriteNull()...)
	data = append(data, amf0WriteString(pkt.StreamName)...)
	data = append(data, amf0WriteNumber(pkt.Start)...)
	data = append(data, amf0WriteNumber(pkt.Duration)...)
	data = append(data, amf0WriteBool(pkt.Reset)...)

	return
}

// GetMessageType .
func (pkt *PlayPacket) GetMessageType() uint8 {
	return RtmpMsgAmf0CommandMessage
}

// GetPreferCsID .
func (pkt *PlayPacket) GetPreferCsID() uint32 {
	return RtmpCidOverStream
}
