package main

import (
	"UtilsTools/identify_panic"
	"log"
	"net"
)

type RtmpConn struct {
	net.Conn
	chunks         map[uint32]*ChunkStream //key csid.
	transactionIds map[float64]string      //key transaction id. value: request command name
	ackWindow      struct {
		ackWindowSize uint32 //
		hasAckedSize  uint64 //size has acked to peer
	}
	recvBytesSum   uint64
	chunkSize      uint32 //default is RTMP_DEFAULT_CHUNK_SIZE. can set by peer.
	objectEncoding float64
}

func NewRtmpSession(c net.Conn) *RtmpConn {
	return &RtmpConn{
		Conn:           c,
		chunks:         make(map[uint32]*ChunkStream),
		chunkSize:      RTMP_DEFAULT_CHUNK_SIZE,
		objectEncoding: RTMP_SIG_AMF0_VER,
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
