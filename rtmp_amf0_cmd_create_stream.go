package main

import (
	"UtilsTools/identify_panic"
	"log"
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

	log.Println("handle amf0 cmd create stream success. comand=", commandName, ", transaction id=", transactionId)

	return
}
