package seal_rtmp_server

import (
	"UtilsTools/identify_panic"
	"log"
	"net"
	"seal/seal_conf"
	"seal/seal_rtmp_server/seal_rtmp_conn"
	"seal/seal_rtmp_server/seal_rtmp_protocol/protocol_stack"
	"sync"
)

type RtmpServer struct {
	Conf *seal_conf.ConfInfoRtmp
	Wg   *sync.WaitGroup
}

func (rtmpServer *RtmpServer) NewRtmpSession(c net.Conn) *seal_rtmp_conn.RtmpConn {
	return &seal_rtmp_conn.RtmpConn{
		Conn:           c,
		TimeOut:        rtmpServer.Conf.TimeOut,
		Chunks:         make(map[uint32]*seal_rtmp_conn.ChunkStream),
		ChunkSize:      protocol_stack.RTMP_DEFAULT_CHUNK_SIZE,
		ObjectEncoding: protocol_stack.RTMP_SIG_AMF0_VER,
		Role:           seal_rtmp_conn.RTMP_ROLE_UNKNOWN,
	}
}

func (rtmpServer *RtmpServer) Start() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}

		rtmpServer.Wg.Done()
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
