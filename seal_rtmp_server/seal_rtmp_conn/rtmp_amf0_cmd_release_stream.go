package seal_rtmp_conn

import (
	"UtilsTools/identify_panic"
	"fmt"
	"log"
	"seal/seal_rtmp_server/seal_rtmp_protocol/amf_serial"
	"seal/seal_rtmp_server/seal_rtmp_protocol/protocol_stack"
	"strings"
)

func (rtmp *RtmpConn) handleAmf0CmdReleaseStream(msg *MessageStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	var offset uint32

	var commandName string
	err, commandName = amf_serial.Amf0ReadString(msg.payload, &offset)

	if err != nil {
		return
	}

	if commandName != protocol_stack.RTMP_AMF0_COMMAND_RELEASE_STREAM {
		fmt.Errorf("handleAmf0CmdReleaseStream, cmd is wrong.cmd=", commandName)
		return
	}

	var transactionId float64
	err, transactionId = amf_serial.Amf0ReadNumber(msg.payload, &offset)
	if err != nil {
		return
	}

	err = amf_serial.Amf0ReadNull(msg.payload, &offset)
	if err != nil {
		return
	}

	var streamNameWithToken string
	err, streamNameWithToken = amf_serial.Amf0ReadString(msg.payload, &offset)
	if err != nil {
		return
	}

	err = rtmp.ResponseAmf0CmdReleaseStream(msg.header.preferCsId, transactionId)
	if err != nil {
		return
	}

	rtmp.Role = RTMP_ROLE_PUBLISH

	streamWithoutToken, tokenStr := ParseStreamName(streamNameWithToken)

	_, ok := MapPublishingStreams.Load(streamWithoutToken)
	if ok {
		err = fmt.Errorf("stream ", streamNameWithToken, " can not publish, becasue has publishing now.")
		return
	} else {
		MapPublishingStreams.Store(streamWithoutToken, tokenStr)
	}

	if err != nil {
		return
	}

	rtmp.StreamInfo.stream = streamWithoutToken
	rtmp.StreamInfo.token = tokenStr

	log.Println("handle amf0 cmd release stream success. new publish role, comand=", commandName, ", transaction id=", transactionId,
		"stream name=", streamWithoutToken, ",token=", tokenStr)

	return
}

func (rtmp *RtmpConn) ResponseAmf0CmdReleaseStream(chunkStreamId uint32, transactionId float64) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	var msg MessageStream

	//msg payload
	msg.payload = append(msg.payload, amf_serial.Amf0WriteString(protocol_stack.RTMP_AMF0_COMMAND_RESULT)...)
	msg.payload = append(msg.payload, amf_serial.Amf0WriteNumber(transactionId)...) //transaction id
	msg.payload = append(msg.payload, amf_serial.Amf0WriteNull()...)
	msg.payload = append(msg.payload, amf_serial.Amf0WriteUndefined()...)

	//msg header
	msg.header.length = uint32(len(msg.payload))
	msg.header.typeId = protocol_stack.RTMP_MSG_AMF0CommandMessage
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

func ParseStreamName(s string) (stream string, token string) {

	const TOKEN_STR = "?token="

	loc := strings.Index(s, TOKEN_STR)
	if loc < 0 {
		stream = s
	} else {
		stream = s[0:loc]
		token = s[loc+len(TOKEN_STR):]
	}

	return
}
