package co

import (
	"log"
	"seal/conf"
	"seal/rtmp/pt"
	"utiltools"
)

func (rc *RtmpConn) ResponseAcknowlegementMsg() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	var pkt pt.AcknowlegementPacket

	pkt.SequenceNumber = uint32(rc.TcpConn.RecvBytesSum)

	err = rc.SendPacket(&pkt, 0, conf.GlobalConfInfo.Rtmp.TimeOut*1000000)
	if err != nil {
		return
	}

	rc.AckWindow.HasAckedSize = rc.TcpConn.RecvBytesSum

	return
}

func (rc *RtmpConn) EstimateNeedSendAcknowlegement() (err error) {

	if rc.AckWindow.AckWindowSize > 0 &&
		((rc.TcpConn.RecvBytesSum - rc.AckWindow.HasAckedSize) > uint64(rc.AckWindow.AckWindowSize)) {
		//response a acknowlegement to peer.
		err = rc.ResponseAcknowlegementMsg()
		if err != nil {
			log.Println("response acknowlegement msg failed to peer.")
			return
		}
	}

	return
}
