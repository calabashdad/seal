package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/protocol"
)

func (rc *RtmpConn) Connect() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	var connectPkg *protocol.ConnectPacket
	var pkt protocol.Packet
	pkt = connectPkg
	err = rc.ExpectMsg(&pkt)
	if err != nil {
		log.Println("expect connect packet failed. err=", err)
		return
	}

	connectPkg = pkt.(*protocol.ConnectPacket)

	log.Println("expect connect pkt success.", connectPkg)

	return
}
