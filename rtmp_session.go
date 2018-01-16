package main

import (
	"UtilsTools/identify_panic"
	"log"
)

func HandleRtmpSession(rtmpSession *RtmpSession) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}

		rtmpSession.Conn.Close()
		log.Println("One RtmpSession finished.remote=", rtmpSession.Conn.RemoteAddr())
	}()

	log.Println("One RtmpSession come in. remote=", rtmpSession.Conn.RemoteAddr())

	err := rtmpSession.HandShake()
	if err != nil {
		log.Println("rtmp handshake failed, err=", err)
		return
	}

	log.Println("rtmp handshake success.remote=", rtmpSession.RemoteAddr())

	err = rtmpSession.Connect()
	if err != nil {
		log.Println("rtmp connect failed, err=", err)
		return
	}

	log.Println("rtmp connect success. remote=", rtmpSession.RemoteAddr())
}
