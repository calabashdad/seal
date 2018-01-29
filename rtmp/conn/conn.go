package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/kernel"
	"seal/rtmp/protocol"
)

type AckWindowSizeS struct {
	AckWindowSize uint32
	HasAckedSize  uint64
}

type RtmpConn struct {
	TcpConn      *kernel.TcpSock
	ChunkStreams map[uint32]*protocol.ChunkStream //key:cs id
	InChunkSize  uint32                           //default 128, set by peer
	OutChunkSize uint32                           //default 128, set by config file.
	Pool         *kernel.MemPool
	AckWindow    AckWindowSizeS
	Requests     map[float64]string //key: transactin id, value:command name
}

func (rc *RtmpConn) Cycle() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	var err error

	err = rc.HandShake()
	if err != nil {
		log.Println("rtmp handshake failed.err=", err)
		return
	}
	log.Println("rtmp handshake success.")

	err = rc.Connect()
	if err != nil {
		log.Println("connect failed. err=", err)
		return
	}
	log.Println("connect success.")

	rc.Loop()

	log.Println("rtmp conn finished, remote=", rc.TcpConn.Conn.RemoteAddr())
}
