package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"reflect"
	"seal/rtmp/pt"
)

//expect msg type, ignore other msgs until the type special has recv success.
func (rc *RtmpConn) ExpectMsg(pkt *pt.Packet) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	for {
		var msg *pt.Message

		err = rc.RecvMsg(&msg)
		if err != nil {
			break
		}

		var pktLocal pt.Packet
		err = rc.DecodeMsg(&msg, &pktLocal)
		if err != nil {
			break
		}

		if reflect.TypeOf(*pkt) == reflect.TypeOf(pktLocal) {
			*pkt = pktLocal
			break
		}

	}

	if err != nil {
		return
	}

	return
}
