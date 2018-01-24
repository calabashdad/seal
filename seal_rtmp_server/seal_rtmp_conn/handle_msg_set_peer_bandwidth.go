package seal_rtmp_conn

import (
	"UtilsTools/identify_panic"
	"log"
)

func (rtmp *RtmpConn) handleSetPeerBandWidth(msg *MessageStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	log.Println("set bandwidth")

	return
}
