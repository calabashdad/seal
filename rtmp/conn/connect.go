package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/pt"
)

func (rc *RtmpConn) Connect() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	var connectPkg *pt.ConnectPacket
	var pkt pt.Packet
	pkt = connectPkg
	err = rc.ExpectMsg(&pkt)
	if err != nil {
		log.Println("expect connect packet failed. err=", err)
		return
	}

	connectPkg = pkt.(*pt.ConnectPacket)

	//todo.
	//parse the params and analysis them.

	log.Println("expect connect pkt success.", connectPkg)

	return
}
