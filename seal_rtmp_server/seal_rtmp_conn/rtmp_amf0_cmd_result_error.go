package seal_rtmp_conn

import (
	"UtilsTools/identify_panic"
	"fmt"
	"log"
	"seal/seal_rtmp_server/seal_rtmp_protocol/amf_serial"
	"seal/seal_rtmp_server/seal_rtmp_protocol/protocol_stack"
)

func (rtmp *RtmpConn) handleAMF0CmdResultError(msg *MessageStream) (err error) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	var offset uint32

	var transactionId float64
	err, transactionId = amf_serial.Amf0ReadNumber(msg.payload, &offset)
	if err != nil {
		return
	}

	requestCommand, ok := rtmp.TransactionIds[transactionId]
	if !ok {
		err = fmt.Errorf("handleAMF0CmdResultError can not find the transaction id.")
		return
	}

	switch requestCommand {
	case protocol_stack.RTMP_AMF0_COMMAND_CONNECT:
		//todo
	case protocol_stack.RTMP_AMF0_COMMAND_CREATE_STREAM:
		//todo
	case protocol_stack.RTMP_AMF0_COMMAND_RELEASE_STREAM, protocol_stack.RTMP_AMF0_COMMAND_FC_PUBLISH, protocol_stack.RTMP_AMF0_COMMAND_UNPUBLISH:
		//todo
	case protocol_stack.RTMP_AMF0_COMMAND_ENABLEVIDEO:
		//todo
	case protocol_stack.RTMP_AMF0_COMMAND_INSERT_KEYFREAME:
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
