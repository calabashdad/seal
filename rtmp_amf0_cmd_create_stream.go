package main

import (
	"UtilsTools/identify_panic"
	"log"
	"math"
)

func (rtmp *RtmpConn) handleAmf0CmdCreateStream(msg *MessageStream) (err error) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	var offset uint32

	var commandName string
	err, commandName = Amf0ReadString(msg.payload, &offset)

	var transactionId float64
	err, transactionId = Amf0ReadNumber(msg.payload, &offset)
	if err != nil {
		return
	}

	err = Amf0ReadNull(msg.payload, &offset)
	if err != nil {
		return
	}

	err = rtmp.ResponseAmf0CmdCreateStream(msg.header.preferCsId, msg.header.streamId, transactionId)
	if err != nil {
		return
	}

	log.Println("handle amf0 cmd create stream success. comand=", commandName, ", transaction id=", transactionId)

	return
}

func (rtmp *RtmpConn) ResponseAmf0CmdCreateStream(chunkStreamId uint32, streamId uint32, transactionId float64) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	var msg MessageStream

	//msg payload
	msg.payload = append(msg.payload, Amf0WriteString(RTMP_AMF0_COMMAND_RESULT)...)
	msg.payload = append(msg.payload, Amf0WriteNumber(transactionId)...) //transaction id
	msg.payload = append(msg.payload, Amf0WriteNull()...)
	msg.payload = append(msg.payload, Amf0WriteNumber(math.Float64frombits(uint64(streamId)))...)

	//msg header
	msg.header.length = uint32(len(msg.payload))
	msg.header.typeId = RTMP_MSG_AMF0CommandMessage
	if chunkStreamId < 2 {
		msg.header.preferCsId = RTMP_CID_ProtocolControl
	} else {
		msg.header.preferCsId = chunkStreamId
	}

	err = rtmp.SendMsg(&msg)
	if err != nil {
		return
	}

	return
}
