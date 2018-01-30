package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/protocol"
)

//OnBwDone is the response of bandwidth done packet
func (rc *RtmpConn) OnBwDone() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	var pkt protocol.OnBwDonePacket

	pkt.CommandName = protocol.RTMP_AMF0_COMMAND_ON_BW_DONE
	pkt.TransactionId = 0

	err = rc.SendPacket(&pkt, 0)
	if err != nil {
		return
	}

	return
}
