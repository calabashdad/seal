package main

import (
	"UtilsTools/identify_panic"
	"log"
	"net"
	"seal/conf"
	"seal/kernel"
	"seal/rtmp/conn"
	"seal/rtmp/protocol"
)

type RtmpServer struct {
}

func (rtmp_server *RtmpServer) Start() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}

		gWGServers.Done()
	}()

	listener, err := net.Listen("tcp", ":"+conf.GlobalConfInfo.Rtmp.Listen)
	if err != nil {
		log.Println("start listen at "+conf.GlobalConfInfo.Rtmp.Listen+" failed. err=", err)
		return
	}
	log.Println("rtmp server start liste at :" + conf.GlobalConfInfo.Rtmp.Listen)

	for {
		net_conn, err := listener.Accept()
		if err != nil {
			log.Println("rtmp server, listen accept failed, err=", err)
			break
		}

		log.Println("one rtmp connection come in, remote=", net_conn.RemoteAddr())

		rtmp_conn := rtmp_server.NewRtmpConnection(net_conn)

		go rtmp_conn.Cycle()
	}
}

func (rtmp_server *RtmpServer) NewRtmpConnection(c net.Conn) *conn.RtmpConn {
	return &conn.RtmpConn{
		TcpConn: &kernel.TcpSock{
			Conn:    c,
			TimeOut: conf.GlobalConfInfo.Rtmp.TimeOut,
		},
		In_chunk_size:  protocol.RTMP_DEFAULT_CHUNK_SIZE,
		Out_chunk_size: protocol.RTMP_DEFAULT_CHUNK_SIZE,
		Pool:           kernel.NewMemPool(),
		Ack_window: conn.AckWindowSize{
			Ack_window_size: 250000,
		},
		Requests: make(map[float64]string),
	}
}