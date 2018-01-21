package main

import (
	"UtilsTools/identify_panic"
	"fmt"
	"log"
)

func (rtmp *RtmpConn) handleAMFCommandAndDataMessage(msg *MessageStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}
	}()

	log.Println("amf0/3 command or amf0/3 data")

	var offset uint32

	// skip 1bytes to decode the amf3 command.
	if RTMP_MSG_AMF3CommandMessage == msg.header.typeId && msg.header.length >= 1 {
		offset += 1
	}

	//read the command name.
	err, command := Amf0ReadString(msg.payload, &offset)
	if err != nil {
		return
	}

	if 0 == len(command) {
		err = fmt.Errorf("Amf0ReadString failed, command is nil.")
		return
	}

	log.Println("msg typeid=", msg.header.typeId, ",command =", command)

	switch command {
	case RTMP_AMF0_COMMAND_RESULT, RTMP_AMF0_COMMAND_ERROR:
		err = rtmp.handleAMF0CmdResultError(msg)
	case RTMP_AMF0_COMMAND_CONNECT:
		err = rtmp.handleAMF0CmdConnect(msg)
	case RTMP_AMF0_COMMAND_CREATE_STREAM:
		err = rtmp.handleAmf0CmdCreateStream(msg)
	case RTMP_AMF0_COMMAND_PLAY:
		//todo.
	case RTMP_AMF0_COMMAND_PAUSE:
		//todo.
	case RTMP_AMF0_COMMAND_RELEASE_STREAM:
		err = rtmp.handleAmf0CmdReleaseStream(msg)
	case RTMP_AMF0_COMMAND_FC_PUBLISH:
		err = rtmp.handleAmf0CmdFcPublish(msg)
	case RTMP_AMF0_COMMAND_PUBLISH:
		//todo.
	case RTMP_AMF0_COMMAND_UNPUBLISH:
		//todo.
	case RTMP_AMF0_COMMAND_KEEPLIVE:
		//todo.
	case RTMP_AMF0_COMMAND_ENABLEVIDEO:
		//todo.
	case RTMP_AMF0_DATA_SET_DATAFRAME, RTMP_AMF0_DATA_ON_METADATA:
		//todo.
	case RTMP_AMF0_DATA_ON_CUSTOMDATA:
		//todo.
	case RTMP_AMF0_COMMAND_CLOSE_STREAM:
		//todo.
	case RTMP_AMF0_COMMAND_ON_BW_DONE:
		//todo
	case RTMP_AMF0_COMMAND_ON_STATUS:
		//todo
	case RTMP_AMF0_COMMAND_INSERT_KEYFREAME:
		//todo
	case RTMP_AMF0_DATA_SAMPLE_ACCESS:
		//todo.
	default:
		log.Println("handleAMFCommandAndDataMessage:unknown command name.", command)
	}

	if err != nil {
		return
	}

	return
}
