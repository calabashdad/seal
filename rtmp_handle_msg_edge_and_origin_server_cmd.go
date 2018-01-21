package main

import (
	"UtilsTools/identify_panic"
	"log"
)

func (rtmp *RtmpConn) handleEdgeAndOriginServerCommand(msg *MessageStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	log.Println("edge and origin server command, remote=", rtmp.Conn.RemoteAddr())

	return
}
