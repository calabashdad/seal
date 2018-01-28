package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/kernel"
	"seal/rtmp/protocol"
)

type AckWindowSize struct {
	Ack_window_size uint32
	Has_acked_size  uint64
}

type RtmpConn struct {
	TcpConn        *kernel.TcpSock
	chunk_streams  map[uint32]*protocol.ChunkStream //key:cs id
	In_chunk_size  uint32                           //default 128, set by peer
	Out_chunk_size uint32                           //default 128, set by config file.
	Pool           *kernel.MemPool
	Ack_window     AckWindowSize
}

func (rtmp_conn *RtmpConn) Cycle() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	var err error

	err = rtmp_conn.HandShake()
	if err != nil {
		log.Println("rtmp handshake failed.err=", err)
		return
	}
	log.Println("rtmp handshake success.")

	log.Println("rtmp conn finished, remote=", rtmp_conn.TcpConn.Conn.RemoteAddr())
}
