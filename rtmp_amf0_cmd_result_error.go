package main

import (
	"UtilsTools/identify_panic"
	"fmt"
	"log"
)

func (rtmp *RtmpConn) handleAMF0CmdResultError(msg *MessageStream) (err error) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	var offset uint32

	var transactionId float64
	err, transactionId = Amf0ReadNumber(msg.payload, &offset)
	if err != nil {
		return
	}

	requestCommand, ok := rtmp.transactionIds[transactionId]
	if !ok {
		err = fmt.Errorf("handleAMF0CmdResultError can not find the transaction id.")
		return
	}

	switch requestCommand {
	case RTMP_AMF0_COMMAND_CONNECT:
		//todo
	case RTMP_AMF0_COMMAND_CREATE_STREAM:
		//todo
	case RTMP_AMF0_COMMAND_RELEASE_STREAM, RTMP_AMF0_COMMAND_FC_PUBLISH, RTMP_AMF0_COMMAND_UNPUBLISH:
		//todo
	case RTMP_AMF0_COMMAND_ENABLEVIDEO:
		//todo
	case RTMP_AMF0_COMMAND_INSERT_KEYFREAME:
		//todo
	default:
		err = fmt.Errorf("handleAMF0CmdResultError, unknown request command name,", requestCommand)
	}

	if err != nil {
		return
	}

	log.Println("handle amf0 result/error success.")

	return
}
