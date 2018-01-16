package main

import (
	"UtilsTools/identify_panic"
	"log"
	"net"
)

type RtmpSession struct {
	net.Conn
	chunks map[uint32]ChunkStruct
}

func NewRtmpSession(c net.Conn) *RtmpSession {
	return &RtmpSession{
		Conn:   c,
		chunks: make(map[uint32]ChunkStruct),
	}
}

func StartRtmpServer() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}

		g_wg.Done()
	}()

	listener, err := net.Listen("tcp", ":"+g_conf_info.Rtmp.Listen)
	if err != nil {
		log.Println("start listen at "+g_conf_info.Rtmp.Listen+" failed. err=", err)
		return
	}

	log.Println("rtmp server start liste at :" + g_conf_info.Rtmp.Listen)

	for {
		netconn, err := listener.Accept()
		if err != nil {
			log.Println("rtmp server, listen accept failed, err=", err)
			break
		}

		rtmpConn := NewRtmpSession(netconn)

		go HandleRtmpSession(rtmpConn)
	}
}
