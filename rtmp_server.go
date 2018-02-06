package main

import (
	"log"
	"net"
	"seal/conf"
	"seal/kernel"
	"seal/rtmp/co"
	"seal/rtmp/pt"
	"utiltools"
)

type RtmpServer struct {
}

func (rs *RtmpServer) Start() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}

		gWgServers.Done()
	}()

	listener, err := net.Listen("tcp", ":"+conf.GlobalConfInfo.Rtmp.Listen)
	if err != nil {
		log.Println("start listen at "+conf.GlobalConfInfo.Rtmp.Listen+" failed. err=", err)
		return
	}
	log.Println("rtmp server start liste at :" + conf.GlobalConfInfo.Rtmp.Listen)

	for {
		netConn, err := listener.Accept()
		if err != nil {
			log.Println("rtmp server, listen accept failed, err=", err)
			break
		}

		log.Println("one rtmp connection come in, remote=", netConn.RemoteAddr())

		rtmpConn := rs.NewRtmpConnection(netConn)

		go rtmpConn.Cycle()
	}
}

func (rtmp_server *RtmpServer) NewRtmpConnection(c net.Conn) *co.RtmpConn {
	return &co.RtmpConn{
		TcpConn: &kernel.TcpSock{
			Conn:    c,
			TimeOut: conf.GlobalConfInfo.Rtmp.TimeOut,
		},
		ChunkStreams: make(map[uint32]*pt.ChunkStream),
		InChunkSize:  pt.RTMP_DEFAULT_CHUNK_SIZE,
		OutChunkSize: pt.RTMP_DEFAULT_CHUNK_SIZE,
		AckWindow: co.AckWindowSizeS{
			AckWindowSize: 250000,
		},
		CmdRequests:     make(map[float64]string),
		Role:            co.RtmpRoleUnknown,
		DefaultStreamId: 1.0,
		ConnectInfo: &co.ConnectInfoS{
			ObjectEncoding: pt.RTMP_SIG_AMF0_VER,
		},
	}
}
