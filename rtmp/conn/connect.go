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

	var connect_pkg protocol.ConnectPacket

	err = rc.ExpectMsg(&connect_pkg)
	if err != nil {
		log.Println("expect connect packet failed. err=", err)
		return
	}
	log.Println("expect connect pkg success.")

	return
}
