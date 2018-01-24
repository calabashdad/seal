package seal_rtmp_conn

import (
	"UtilsTools/identify_panic"
	"log"
)

func (rtmp *RtmpConn) handleAbortMsg(msg *MessageStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	log.Println("abort msg, remote=", rtmp.Conn.RemoteAddr())

	return
}
