package main

import (
	"UtilsTools/identify_panic"
	"log"
)

func (rtmp *RtmpConn) handleAcknowlegement(msg *MessageStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	log.Println("acknowlegement")

	return
}
