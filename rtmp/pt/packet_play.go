package pt

import (
	"fmt"
)

type PlayPacket struct {
	/**
	 * Name of the command. Set to “play”.
	 */
	Command_name string
	/**
	 * Transaction ID set to 0.
	 */
	Transaction_id float64
	/**
	 * Command information does not exist. Set to null type.
	 */
	Command_object Amf0Object // null
	/**
	 * Name of the stream to play.
	 * To play video (FLV) files, specify the name of the stream without a file
	 *       extension (for example, "sample").
	 * To play back MP3 or ID3 tags, you must precede the stream name with mp3:
	 *       (for example, "mp3:sample".)
	 * To play H.264/AAC files, you must precede the stream name with mp4: and specify the
	 *       file extension. For example, to play the file sample.m4v, specify
	 *       "mp4:sample.m4v"
	 */
	StreamName string
	/**
	 * An optional parameter that specifies the start time in seconds.
	 * The default value is -2, which means the subscriber first tries to play the live
	 *       stream specified in the Stream Name field. If a live stream of that name is
	 *       not found, it plays the recorded stream specified in the Stream Name field.
	 * If you pass -1 in the Start field, only the live stream specified in the Stream
	 *       Name field is played.
	 * If you pass 0 or a positive number in the Start field, a recorded stream specified
	 *       in the Stream Name field is played beginning from the time specified in the
	 *       Start field.
	 * If no recorded stream is found, the next item in the playlist is played.
	 */
	Start float64
	/**
	 * An optional parameter that specifies the duration of playback in seconds.
	 * The default value is -1. The -1 value means a live stream is played until it is no
	 *       longer available or a recorded stream is played until it ends.
	 * If u pass 0, it plays the single frame since the time specified in the Start field
	 *       from the beginning of a recorded stream. It is assumed that the value specified
	 *       in the Start field is equal to or greater than 0.
	 * If you pass a positive number, it plays a live stream for the time period specified
	 *       in the Duration field. After that it becomes available or plays a recorded
	 *       stream for the time specified in the Duration field. (If a stream ends before the
	 *       time specified in the Duration field, playback ends when the stream ends.)
	 * If you pass a negative number other than -1 in the Duration field, it interprets the
	 *       value as if it were -1.
	 */
	Duration float64
	/**
	 * An optional Boolean value or number that specifies whether to flush any
	 * previous playlist.
	 */
	Reset bool
}

func (pkt *PlayPacket) Decode(data []uint8) (err error) {
	var maxOffset uint32
	maxOffset = uint32(len(data)) - 1

	var offset uint32

	err, pkt.Command_name = Amf0ReadString(data, &offset)
	if err != nil {
		return
	}

	if RTMP_AMF0_COMMAND_PLAY != pkt.Command_name {
		err = fmt.Errorf("decode play packet, command name is not play.actully=%s", pkt.Command_name)
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

	err, pkt.StreamName = Amf0ReadString(data, &offset)
	if err != nil {
		return
	}

	if maxOffset-offset > (1 + 8) { // number need at least 1(marker) + 8(number)
		err, pkt.Start = Amf0ReadNumber(data, &offset)
		if err != nil {
			return
		}
	}

	if maxOffset-offset > (1 + 8) {
		err, pkt.Duration = Amf0ReadNumber(data, &offset)
		if err != nil {
			return
		}
	}

	if offset >= uint32(len(data)) {
		return
	}

	if maxOffset-offset >= 2 { //because the bool type need 2 bytes at least

		var v interface{}
		var marker uint8
		err, v = Amf0ReadAny(data, &marker, &offset)
		if err != nil {
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
func (pkt *PlayPacket) Encode() (data []uint8) {

	data = append(data, Amf0WriteString(pkt.Command_name)...)
	data = append(data, Amf0WriteNumber(pkt.Transaction_id)...)
	data = append(data, Amf0WriteNull()...)
	data = append(data, Amf0WriteString(pkt.StreamName)...)
	data = append(data, Amf0WriteNumber(pkt.Start)...)
	data = append(data, Amf0WriteNumber(pkt.Duration)...)
	data = append(data, Amf0WriteBool(pkt.Reset)...)

	return
}
func (pkt *PlayPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}
func (pkt *PlayPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverStream
}
