package main

import (
	"UtilsTools/identify_panic"
	"log"
)

func HandleRtmpConn(conn *RtmpConn) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}

		conn.Conn.Close()
		log.Println("One RtmpConn finished.remote=", conn.Conn.RemoteAddr())
	}()

	log.Println("One RtmpConn come in. remote=", conn.Conn.RemoteAddr())

	err := conn.HandShake()
	if err != nil {
		log.Println("rtmp handshake failed, err=", err)
		return
	}

	log.Println("rtmp handshake success.remote=", conn.RemoteAddr())
}
