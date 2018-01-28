package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/protocol"
)

func (rc *RtmpConn) OnRecvMsg(msg *protocol.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	if rc.Ack_window.Ack_window_size > 0 &&
		((rc.TcpConn.RecvBytesSum - rc.Ack_window.Has_acked_size) > uint64(rc.Ack_window.Ack_window_size)) {
		//response a acknowlegement to peer.
		err = rc.ResponseAcknowlegementMsg()
		if err != nil {
			log.Println("response acknowlegement msg failed to peer.")
			return
		}
	}

	return
}
