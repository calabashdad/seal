package main

import (
	"UtilsTools/identify_panic"
	"log"
	"net"
	"seal/conf"
	"seal/seal_rtmp_server/seal_rtmp_conn"
	"seal/seal_rtmp_server/seal_rtmp_protocol/protocol_stack"
)

type RtmpServer struct {
	Conf *conf.ConfInfoRtmp
}

func (rtmpServer *RtmpServer) NewRtmpSession(c net.Conn) *seal_rtmp_conn.RtmpConn {
	return &seal_rtmp_conn.RtmpConn{
		Conn:          c,
		TimeOut:       rtmpServer.Conf.TimeOut,
		Chunks:        make(map[uint32]*seal_rtmp_conn.ChunkStream),
		ChunkSize:     protocol_stack.RTMP_CHUNKSIZE_MIN,
		ChunkSizeConf: rtmpServer.Conf.ChunkSize,
		Role:          seal_rtmp_conn.RTMP_ROLE_UNKNOWN,
	}
}

func (rtmpServer *RtmpServer) Start() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}

		g_wg.Done()
	}()

	listener, err := net.Listen("tcp", ":"+rtmpServer.Conf.Listen)
	if err != nil {
		log.Println("start listen at "+rtmpServer.Conf.Listen+" failed. err=", err)
		return
	}

	log.Println("rtmp server start liste at :" + rtmpServer.Conf.Listen)

	for {
		netConn, err := listener.Accept()
		if err != nil {
			log.Println("rtmp server, listen accept failed, err=", err)
			break
		}

		rtmpConn := rtmpServer.NewRtmpSession(netConn)

		go rtmpConn.HandleRtmpSession()
	}
}
