package main

import (
	"UtilsTools/identify_panic"
	"log"
)

func (rtmpSession *RtmpConn) RtmpMsgLoop() (err error) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ", panic at ", identify_panic.IdentifyPanic())
		}
	}()

	for {
		var chunkStreamId uint32
		err, chunkStreamId = rtmpSession.RecvMsg()
		if err != nil {
			break
		}

		err = rtmpSession.DecodeAndHanleMsg(chunkStreamId)
		if err != nil {
			break
		}
	}

	if err != nil {
		return
	}

	return
}
