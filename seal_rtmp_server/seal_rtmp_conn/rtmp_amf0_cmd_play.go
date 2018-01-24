package seal_rtmp_conn

import (
	"UtilsTools/identify_panic"
	"fmt"
	"log"
	"seal/seal_rtmp_server/seal_rtmp_protocol/amf_serial"
	"seal/seal_rtmp_server/seal_rtmp_protocol/protocol_stack"
)

func (rtmp *RtmpConn) handleAmf0CmdPlay(msg *MessageStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	var offset uint32
	maxOffset := uint32(len(msg.payload)) - 1

	var commandName string
	err, commandName = amf_serial.Amf0ReadString(msg.payload, &offset)
	if err != nil {
		return
	}

	if commandName != protocol_stack.RTMP_AMF0_COMMAND_PLAY {
		err = fmt.Errorf("handleAmf0CmdPlay commandName is wrong.", commandName)
		return
	}

	var transactionId float64
	err, transactionId = amf_serial.Amf0ReadNumber(msg.payload, &offset)
	if err != nil {
		return
	}

	amf_serial.Amf0ReadNull(msg.payload, &offset)
	if err != nil {
		return
	}

	var streamNameWithToken string
	err, streamNameWithToken = amf_serial.Amf0ReadString(msg.payload, &offset)
	if err != nil {
		return
	}

	rtmp.StreamInfo.stream, rtmp.StreamInfo.token = ParseStreamName(streamNameWithToken)

	var startTime float64
	if maxOffset-offset > (1 + 8) {
		err, startTime = amf_serial.Amf0ReadNumber(msg.payload, &offset)
		if err != nil {
			return
		}
	}

	var durationOfPlayback float64
	if maxOffset-offset > (1 + 8) {
		err, durationOfPlayback = amf_serial.Amf0ReadNumber(msg.payload, &offset)
		if err != nil {
			return
		}
	}

	log.Println("handle amf0 cmd play success. ",
		",transactionId=", transactionId,
		",streamName=", streamNameWithToken,
		",startTime=", startTime,
		",optional param, durationOfPlayback=", durationOfPlayback,
		",new player come in, remote =", rtmp.Conn.RemoteAddr())

	if maxOffset-offset < 2 { //becase after is bool or number, at least is 2 bytes for bool.
		return
	}

	var resetPlayList bool
	var obj interface{}
	var objMarker uint8
	err, obj = amf_serial.Amf0ReadAny(msg.payload, &objMarker, &offset)
	if err != nil {
		return
	}
	if protocol_stack.RTMP_AMF0_Boolean == objMarker {
		resetPlayList = obj.(bool)
	} else if protocol_stack.RTMP_AMF0_Number == objMarker {
		resetPlayList = 0 != obj.(float64)
	} else {
		err = fmt.Errorf("handleAmf0CmdPlay, reset value type is error.", objMarker)
		return
	}

	if err != nil {
		return
	}

	log.Println("amf0 play, resetPlayList=", resetPlayList)

	return
}
