package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/protocol"
)

func (rc *RtmpConn) ResponseAcknowlegementMsg() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	var pkt protocol.AcknowlegementPacket

	pkt.Sequence_number = uint32(rc.TcpConn.RecvBytesSum)

	err = rc.SendPacket(&pkt, 0)
	if err != nil {
		return
	}

	rc.Ack_window.Has_acked_size = rc.TcpConn.RecvBytesSum

	return
}
