package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"reflect"
	"seal/rtmp/protocol"
)

//expect msg type, ignore other msgs until the type special has recv success.
func (rc *RtmpConn) ExpectMsg(pkt protocol.Packet) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	for {
		var msg *protocol.Message

		err = rc.RecvMsg(&msg)
		if err != nil {
			break
		}

		var pkt_local protocol.Packet
		err = rc.DecodeMsg(&msg, &pkt_local)
		if err != nil {
			break
		}

		if reflect.TypeOf(pkt_local) == reflect.TypeOf(pkt) {
			pkt = pkt_local
			break
		}

	}

	if err != nil {
		return
	}

	return
}
