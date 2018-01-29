package conn

import (
	"seal/rtmp/protocol"
	"UtilsTools/identify_panic"
	"log"
)

func (rc *RtmpConn) OnBwDone() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	var pkt protocol.OnBwDonePacket
	
	return
}
