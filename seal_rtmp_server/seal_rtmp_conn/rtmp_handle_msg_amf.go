package seal_rtmp_conn

import (
	"UtilsTools/identify_panic"
	"fmt"
	"log"
	"seal/seal_rtmp_server/seal_rtmp_protocol/amf_serial"
	"seal/seal_rtmp_server/seal_rtmp_protocol/protocol_stack"
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
	if protocol_stack.RTMP_MSG_AMF3CommandMessage == msg.header.typeId && msg.header.length >= 1 {
		offset += 1
	}

	//read the command name.
	err, command := amf_serial.Amf0ReadString(msg.payload, &offset)
	if err != nil {
		return
	}

	if 0 == len(command) {
		err = fmt.Errorf("Amf0ReadString failed, command is nil.")
		return
	}

	log.Println("msg typeid=", msg.header.typeId, ",command =", command)

	switch command {
	case protocol_stack.RTMP_AMF0_COMMAND_RESULT, protocol_stack.RTMP_AMF0_COMMAND_ERROR:
		err = rtmp.handleAMF0CmdResultError(msg)
	case protocol_stack.RTMP_AMF0_COMMAND_CONNECT:
		err = rtmp.handleAMF0CmdConnect(msg)
	case protocol_stack.RTMP_AMF0_COMMAND_CREATE_STREAM:
		err = rtmp.handleAmf0CmdCreateStream(msg)
	case protocol_stack.RTMP_AMF0_COMMAND_PLAY:
		//todo.
	case protocol_stack.RTMP_AMF0_COMMAND_PAUSE:
		//todo.
	case protocol_stack.RTMP_AMF0_COMMAND_RELEASE_STREAM:
		err = rtmp.handleAmf0CmdReleaseStream(msg)
	case protocol_stack.RTMP_AMF0_COMMAND_FC_PUBLISH:
		err = rtmp.handleAmf0CmdFcPublish(msg)
	case protocol_stack.RTMP_AMF0_COMMAND_PUBLISH:
		err = rtmp.handleAmf0CmdPublish(msg)
	case protocol_stack.RTMP_AMF0_COMMAND_UNPUBLISH:
		//todo.
	case protocol_stack.RTMP_AMF0_COMMAND_KEEPLIVE:
		//todo.
	case protocol_stack.RTMP_AMF0_COMMAND_ENABLEVIDEO:
		//todo.
	case protocol_stack.RTMP_AMF0_DATA_SET_DATAFRAME, protocol_stack.RTMP_AMF0_DATA_ON_METADATA:
		//todo.
	case protocol_stack.RTMP_AMF0_DATA_ON_CUSTOMDATA:
		//todo.
	case protocol_stack.RTMP_AMF0_COMMAND_CLOSE_STREAM:
		//todo.
	case protocol_stack.RTMP_AMF0_COMMAND_ON_BW_DONE:
		//todo
	case protocol_stack.RTMP_AMF0_COMMAND_ON_STATUS:
		//todo
	case protocol_stack.RTMP_AMF0_COMMAND_INSERT_KEYFREAME:
		//todo
	case protocol_stack.RTMP_AMF0_DATA_SAMPLE_ACCESS:
		//todo.
	default:
		log.Println("handleAMFCommandAndDataMessage:unknown command name.", command)
	}

	if err != nil {
		return
	}

	return
}
