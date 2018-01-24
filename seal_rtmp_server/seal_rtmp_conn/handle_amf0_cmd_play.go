package seal_rtmp_conn

import (
	"UtilsTools/identify_panic"
	"encoding/binary"
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

	rtmp.StreamInfo.stream, rtmp.StreamInfo.token = handleParseStreamName(streamNameWithToken)

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

	rtmp.Role = RTMP_ROLE_PALY
	rtmp.msgChan = make(chan *MessageStream, 256)
	registerRes := rtmp.PlayerRegistePublishStream()
	if !registerRes {
		err = fmt.Errorf("can not play, because the stream is not publishing.", rtmp.StreamInfo.stream)
		return
	}

	log.Println("handle amf0 cmd play success. ",
		",transactionId=", transactionId,
		",streamName=", streamNameWithToken,
		",startTime=", startTime,
		",optional param, durationOfPlayback=", durationOfPlayback,
		",new player come in, remote =", rtmp.Conn.RemoteAddr())

	if maxOffset-offset >= 2 { //becase after is bool or number, at least is 2 bytes for bool.
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
	}

	err = rtmp.ResponseAmf0CmdPlay(msg)
	if err != nil {
		return
	}

	log.Println("response to play client success. now playing loop....")

	err = rtmp.handlePlayLoop()
	if err != nil {
		log.Println("player msg loop quit ,stream=", rtmp.StreamInfo.stream, ",err=", err)
		return
	}

	return
}

func (rtmp *RtmpConn) ResponseAmf0CmdPlay(msg *MessageStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	err = rtmp.ResponseAmf0CmdPlayStreamBegin(msg.header.preferCsId, msg.header.streamId)
	if err != nil {
		return
	}

	err = rtmp.ResponseAmf0CmdPlayOnStatusNetStreamPlayReset(msg.header.preferCsId, msg.header.streamId)
	if err != nil {
		return
	}

	err = rtmp.ResponseAmf0CmdPlayOnStatusNetStreamPlayStart(msg.header.preferCsId, msg.header.streamId)
	if err != nil {
		return
	}

	err = rtmp.ResponseAmf0CmdPlayRtmpSampleAccess(msg.header.preferCsId, msg.header.streamId)
	if err != nil {
		return
	}

	err = rtmp.ResponseAmf0CmdPlayOnStatusNetStreamDataStart(msg.header.preferCsId, msg.header.streamId)
	if err != nil {
		return
	}

	return
}

func (rtmp *RtmpConn) ResponseAmf0CmdPlayStreamBegin(chunkStreamId uint32, streamId uint32) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	var userCtrl UserControlMsg
	userCtrl.eventType = protocol_stack.SrcPCUCStreamBegin
	userCtrl.eventData = streamId

	var msg MessageStream

	//msg payload
	msg.payload = make([]uint8, 6)

	var offset uint32
	binary.BigEndian.PutUint16(msg.payload[offset:offset+2], userCtrl.eventType)
	offset += 2
	binary.BigEndian.PutUint32(msg.payload[offset:offset+4], userCtrl.eventData)
	offset += 4

	//msg header
	msg.header.length = uint32(len(msg.payload))
	msg.header.typeId = protocol_stack.RTMP_MSG_UserControlMessage
	msg.header.streamId = 0
	if chunkStreamId < 2 {
		msg.header.preferCsId = protocol_stack.RTMP_CID_ProtocolControl
	} else {
		msg.header.preferCsId = chunkStreamId
	}

	err = rtmp.SendMsg(&msg)
	if err != nil {
		return
	}

	return
}

func (rtmp *RtmpConn) ResponseAmf0CmdPlayOnStatusNetStreamPlayReset(chunkStreamId uint32, streamId uint32) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	var msg MessageStream

	//payload
	msg.payload = append(msg.payload, amf_serial.Amf0WriteString(protocol_stack.RTMP_AMF0_COMMAND_ON_STATUS)...)
	msg.payload = append(msg.payload, amf_serial.Amf0WriteNumber(0.0)...) //transaction id, set to 0
	msg.payload = append(msg.payload, amf_serial.Amf0WriteNull()...)

	var objs []amf_serial.Amf0Object

	objs = append(objs, amf_serial.Amf0Object{
		PropertyName: protocol_stack.StatusLevel,
		Value:        protocol_stack.StatusLevelStatus,
		ValueType:    protocol_stack.RTMP_AMF0_String,
	})

	objs = append(objs, amf_serial.Amf0Object{
		PropertyName: protocol_stack.StatusCode,
		Value:        protocol_stack.StatusCodeStreamReset,
		ValueType:    protocol_stack.RTMP_AMF0_String,
	})

	objs = append(objs, amf_serial.Amf0Object{
		PropertyName: protocol_stack.StatusDescription,
		Value:        "Playing and resetting stream.",
		ValueType:    protocol_stack.RTMP_AMF0_String,
	})
	objs = append(objs, amf_serial.Amf0Object{
		PropertyName: protocol_stack.StatusDetails,
		Value:        "stream",
		ValueType:    protocol_stack.RTMP_AMF0_String,
	})
	objs = append(objs, amf_serial.Amf0Object{
		PropertyName: protocol_stack.StatusClientId,
		Value:        protocol_stack.RTMP_SIG_CLIENT_ID,
		ValueType:    protocol_stack.RTMP_AMF0_String,
	})

	msg.payload = append(msg.payload, amf_serial.Amf0WriteObject(objs)...)

	//header
	msg.header.length = uint32(len(msg.payload))
	msg.header.typeId = protocol_stack.RTMP_MSG_AMF0CommandMessage
	msg.header.streamId = streamId
	if chunkStreamId < 2 {
		msg.header.preferCsId = protocol_stack.RTMP_CID_ProtocolControl
	} else {
		msg.header.preferCsId = chunkStreamId
	}

	err = rtmp.SendMsg(&msg)
	if err != nil {
		return
	}

	return
}

func (rtmp *RtmpConn) ResponseAmf0CmdPlayOnStatusNetStreamPlayStart(chunkStreamId uint32, streamId uint32) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	var msg MessageStream

	//payload
	msg.payload = append(msg.payload, amf_serial.Amf0WriteString(protocol_stack.RTMP_AMF0_COMMAND_ON_STATUS)...)
	msg.payload = append(msg.payload, amf_serial.Amf0WriteNumber(0.0)...) //transaction id, set to 0
	msg.payload = append(msg.payload, amf_serial.Amf0WriteNull()...)

	var objs []amf_serial.Amf0Object

	objs = append(objs, amf_serial.Amf0Object{
		PropertyName: protocol_stack.StatusLevel,
		Value:        protocol_stack.StatusLevelStatus,
		ValueType:    protocol_stack.RTMP_AMF0_String,
	})

	objs = append(objs, amf_serial.Amf0Object{
		PropertyName: protocol_stack.StatusCode,
		Value:        protocol_stack.StatusCodeStreamStart,
		ValueType:    protocol_stack.RTMP_AMF0_String,
	})

	objs = append(objs, amf_serial.Amf0Object{
		PropertyName: protocol_stack.StatusDescription,
		Value:        "Started playing stream.",
		ValueType:    protocol_stack.RTMP_AMF0_String,
	})

	objs = append(objs, amf_serial.Amf0Object{
		PropertyName: protocol_stack.StatusDetails,
		Value:        "stream",
		ValueType:    protocol_stack.RTMP_AMF0_String,
	})

	objs = append(objs, amf_serial.Amf0Object{
		PropertyName: protocol_stack.StatusClientId,
		Value:        protocol_stack.RTMP_SIG_CLIENT_ID,
		ValueType:    protocol_stack.RTMP_AMF0_String,
	})

	msg.payload = append(msg.payload, amf_serial.Amf0WriteObject(objs)...)

	//header
	msg.header.length = uint32(len(msg.payload))
	msg.header.typeId = protocol_stack.RTMP_MSG_AMF0CommandMessage
	msg.header.streamId = streamId
	if chunkStreamId < 2 {
		msg.header.preferCsId = protocol_stack.RTMP_CID_ProtocolControl
	} else {
		msg.header.preferCsId = chunkStreamId
	}

	err = rtmp.SendMsg(&msg)
	if err != nil {
		return
	}

	return
}

func (rtmp *RtmpConn) ResponseAmf0CmdPlayRtmpSampleAccess(chunkStreamId uint32, streamId uint32) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	var msg MessageStream

	//payload
	msg.payload = append(msg.payload, amf_serial.Amf0WriteString(protocol_stack.RTMP_AMF0_DATA_SAMPLE_ACCESS)...)
	msg.payload = append(msg.payload, amf_serial.Amf0WriteBool(true)...)
	msg.payload = append(msg.payload, amf_serial.Amf0WriteBool(true)...)

	//header
	msg.header.length = uint32(len(msg.payload))
	msg.header.typeId = protocol_stack.RTMP_MSG_AMF0CommandMessage
	msg.header.streamId = streamId
	if chunkStreamId < 2 {
		msg.header.preferCsId = protocol_stack.RTMP_CID_ProtocolControl
	} else {
		msg.header.preferCsId = chunkStreamId
	}

	err = rtmp.SendMsg(&msg)
	if err != nil {
		return
	}
	return
}

func (rtmp *RtmpConn) ResponseAmf0CmdPlayOnStatusNetStreamDataStart(chunkStreamId uint32, streamId uint32) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	var msg MessageStream

	//payload
	msg.payload = append(msg.payload, amf_serial.Amf0WriteString(protocol_stack.RTMP_AMF0_COMMAND_ON_STATUS)...)
	msg.payload = append(msg.payload, amf_serial.Amf0WriteNumber(0.0)...) //transaction id, set to 0
	msg.payload = append(msg.payload, amf_serial.Amf0WriteNull()...)

	var objs []amf_serial.Amf0Object

	objs = append(objs, amf_serial.Amf0Object{
		PropertyName: protocol_stack.StatusLevel,
		Value:        protocol_stack.StatusLevelStatus,
		ValueType:    protocol_stack.RTMP_AMF0_String,
	})

	objs = append(objs, amf_serial.Amf0Object{
		PropertyName: protocol_stack.StatusCode,
		Value:        protocol_stack.StatusCodeDataStart,
		ValueType:    protocol_stack.RTMP_AMF0_String,
	})

	objs = append(objs, amf_serial.Amf0Object{
		PropertyName: protocol_stack.StatusDescription,
		Value:        "Started playing stream data.",
		ValueType:    protocol_stack.RTMP_AMF0_String,
	})

	msg.payload = append(msg.payload, amf_serial.Amf0WriteObject(objs)...)

	//header
	msg.header.length = uint32(len(msg.payload))
	msg.header.typeId = protocol_stack.RTMP_MSG_AMF0CommandMessage
	msg.header.streamId = streamId
	if chunkStreamId < 2 {
		msg.header.preferCsId = protocol_stack.RTMP_CID_ProtocolControl
	} else {
		msg.header.preferCsId = chunkStreamId
	}

	err = rtmp.SendMsg(&msg)
	if err != nil {
		return
	}

	return
}
