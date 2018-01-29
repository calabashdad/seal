package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/protocol"
)

func (rc *RtmpConn) SetWindowAckSize(ackSize uint32) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	var pkt protocol.SetWindowAckSizePacket
	pkt.Ackowledgement_window_size = ackSize

	err = rc.SendPacket(&pkt, 0)
	if err != nil {
		return
	}

	return
}
