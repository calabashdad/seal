package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"reflect"
	"seal/rtmp/protocol"
)

//expect msg type, ignore other msgs until the type special has recv success.
func (rc *RtmpConn) ExpectMsg(pkt *protocol.Packet) (err error) {
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

		var pktLocal protocol.Packet
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
