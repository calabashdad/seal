package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/pt"
)

func (rc *RtmpConn) ResponseAcknowlegementMsg() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	var pkt pt.AcknowlegementPacket

	pkt.Sequence_number = uint32(rc.TcpConn.RecvBytesSum)

	err = rc.SendPacket(&pkt, 0)
	if err != nil {
		return
	}

	rc.AckWindow.HasAckedSize = rc.TcpConn.RecvBytesSum

	return
}
