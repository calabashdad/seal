package main

import (
	"log"
	"net"
	"seal/conf"
	"seal/rtmp/co"

	"github.com/calabashdad/utiltools"
)

type rtmpServer struct {
}

func (rs *rtmpServer) Start() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}

		gGuards.Done()
	}()

	listener, err := net.Listen("tcp", ":"+conf.GlobalConfInfo.Rtmp.Listen)
	if err != nil {
		log.Println("start listen at "+conf.GlobalConfInfo.Rtmp.Listen+" failed. err=", err)
		return
	}
	log.Println("rtmp server start liste at :" + conf.GlobalConfInfo.Rtmp.Listen)

	for {
		if netConn, err := listener.Accept(); err != nil {
			log.Println("rtmp server, listen accept failed, err=", err)
			break
		} else {
			log.Println("one rtmp connection come in, remote=", netConn.RemoteAddr())
			rtmpConn := co.NewRtmpConnection(netConn)
			go rtmpConn.Cycle()
		}
	}

	log.Println("rtmp server quit, err=", err)
}
