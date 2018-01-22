package seal_rtmp_conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/seal_rtmp_server/seal_rtmp_protocol/amf_serial"
	"seal/seal_rtmp_server/seal_rtmp_protocol/protocol_stack"
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

	var transactionId float64
	err, transactionId = amf_serial.Amf0ReadNumber(msg.payload, &offset)
	if err != nil {
		return
	}

	err = amf_serial.Amf0ReadNull(msg.payload, &offset)
	if err != nil {
		return
	}

	var streamName string
	err, streamName = amf_serial.Amf0ReadString(msg.payload, &offset)
	if err != nil {
		return
	}

	err = rtmp.ResponseAmf0CmdReleaseStream(msg.header.preferCsId, transactionId)
	if err != nil {
		return
	}

	rtmp.Role = RTMP_ROLE_PUBLISH

	log.Println("handle amf0 cmd release stream success. publish role, comand=", commandName, ", transaction id=", transactionId,
		"stream name=", streamName)

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
