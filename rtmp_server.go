package main

import (
	"UtilsTools/identify_panic"
	"log"
	"net"
)

type RtmpConn struct {
	net.Conn
}

func NewRtmpConn(c net.Conn) *RtmpConn {
	return &RtmpConn{c}
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

		rtmpConn := NewRtmpConn(netconn)

		go HandleRtmpConn(rtmpConn)
	}
}
