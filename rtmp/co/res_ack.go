package co

import (
	"log"
	"seal/conf"
	"seal/rtmp/pt"

	"github.com/calabashdad/utiltools"
)

func (rc *RtmpConn) ResponseAcknowlegementMsg() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	var pkt pt.AcknowlegementPacket

	pkt.SequenceNumber = uint32(rc.tcpConn.GetRecvBytesSum())

	if err = rc.SendPacket(&pkt, 0, conf.GlobalConfInfo.Rtmp.TimeOut*1000000); err != nil {
		return
	}

	rc.ack.hasAckedSize = rc.tcpConn.GetRecvBytesSum()

	return
}

func (rc *RtmpConn) estimateNeedSendAcknowlegement() (err error) {

	if rc.ack.ackWindowSize > 0 &&
		((rc.tcpConn.GetRecvBytesSum() - rc.ack.hasAckedSize) > uint64(rc.ack.ackWindowSize)) {
		//response a acknowlegement to peer.
		if err = rc.ResponseAcknowlegementMsg(); err != nil {
			log.Println("response acknowlegement msg failed to peer.")
			return
		}
	}

	return
}
