package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/protocol"
)

func (rc *RtmpConn) RecvMsg() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	for {
		var msg *protocol.Message

		err = rc.RecvInterlacedMsg(msg)
		if err != nil {
			log.Println("recv interlance msg faild.")
			return
		}

		if nil == msg {
			continue
		}

		if msg.Size <= 0 || msg.Header.Payload_length <= 0 {
			log.Println("ignore empty msg.")
			continue
		}

		err = rc.OnRecvMsg(msg)
		if err != nil {
			return
		}

	}

	return
}
