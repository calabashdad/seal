package seal_rtmp_conn

import (
	"UtilsTools/identify_panic"
	"log"
)

func (rtmp *RtmpConn) handleMsgVideo(msg *MessageStream) (err error) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	//log.Println("recv video data, timestamp=", msg.header.timestamp)

	return
}
