package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/pt"
)

func (rc *RtmpConn) RecvMsg(msg **pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	for {

		err = rc.RecvInterlacedMsg(msg)
		if err != nil {
			log.Println("recv interlance msg faild.")
			return
		}

		if nil == *msg {
			//has not recv an entire msg.
			continue
		}

		err = rc.OnRecvMsg(msg)
		if err != nil {
			return
		}

		break
	}

	return
}
