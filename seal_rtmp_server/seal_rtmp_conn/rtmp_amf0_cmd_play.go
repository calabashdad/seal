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

	if err != nil {
		return
	}

	log.Println("handle amf0 cmd play success. ", "transactionId=", transactionId,
		",new player come in, remote =", rtmp.Conn.RemoteAddr())

	return
}
