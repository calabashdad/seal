package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/kernel"
)

type RtmpConn struct {
	TcpConn *kernel.TcpSock
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
