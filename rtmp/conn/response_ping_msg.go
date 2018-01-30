package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/pt"
)

func (rc *RtmpConn) ResponsePingMsg(timeStamp uint32) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	var pkt pt.UserControlPacket

	pkt.Event_type = pt.SrcPCUCPingResponse
	pkt.Event_data = timeStamp

	err = rc.SendPacket(&pkt, 0)
	if err != nil {
		return
	}

	return
}
