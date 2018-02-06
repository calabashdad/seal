package co

import (
	"log"
	"seal/rtmp/pt"
	"utiltools"
)

func (rc *RtmpConn) msgAbort(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("MsgAbort")

	if nil == msg {
		return
	}

	return
}
