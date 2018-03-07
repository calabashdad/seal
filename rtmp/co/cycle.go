package co

import (
	"log"
	"net"
	"seal/conf"
	"seal/kernel"
	"seal/rtmp/pt"

	"github.com/calabashdad/utiltools"
)

type ackWindowSize struct {
	ackWindowSize uint32
	hasAckedSize  uint64
}

type connectInfo struct {
	tcURL          string
	pageURL        string
	swfURL         string
	app            string
	objectEncoding float64
}

// RtmpConn rtmp connection info
type RtmpConn struct {
	tcpConn         *kernel.TCPSock
	chunkStreams    map[uint32]*pt.ChunkStream //key:cs id
	inChunkSize     uint32                     //default 128, set by peer
	outChunkSize    uint32                     //default 128, set by config file.
	ack             ackWindowSize
	cmdRequests     map[float64]string //command requests.key: transactin id, value:command name
	role            uint8              //publisher or player.
	streamName      string
	tokenStr        string        //token str for authentication. it's optional.
	duration        float64       //for player.used to specified the stop when exceed the duration.
	defaultStreamID float64       //default stream id for request.
	connInfo        *connectInfo  //connect info.
	source          *SourceStream //data source info.
	consumer        *Consumer     //for consumer, like player.
}

// NewRtmpConnection create rtmp conncetion
func NewRtmpConnection(c net.Conn) *RtmpConn {
	return &RtmpConn{
		tcpConn: &kernel.TCPSock{
			Conn: c,
		},
		chunkStreams: make(map[uint32]*pt.ChunkStream),
		inChunkSize:  pt.RtmpDefalutChunkSize,
		outChunkSize: pt.RtmpDefalutChunkSize,
		ack: ackWindowSize{
			ackWindowSize: 250000,
		},
		cmdRequests:     make(map[float64]string),
		role:            pt.RtmpRoleUnknown,
		defaultStreamID: 1.0,
		connInfo: &connectInfo{
			objectEncoding: pt.RtmpSigAmf0Ver,
		},
	}
}

// Cycle rtmp service cycle
func (rc *RtmpConn) Cycle() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	rc.tcpConn.SetRecvTimeout(conf.GlobalConfInfo.Rtmp.TimeOut * 1000 * 1000)
	rc.tcpConn.SetSendTimeout(conf.GlobalConfInfo.Rtmp.TimeOut * 1000 * 1000)

	var err error

	if err = rc.handShake(); err != nil {
		log.Println("rtmp handshake failed.err=", err)
		return
	}
	log.Println("rtmp handshake success.")

	for {
		// notice that the payload has not alloced at init.
		// one msg alloc once, and do not copy to improve performance.
		msg := &pt.Message{}

		if err = rc.recvMsg(&msg.Header, &msg.Payload); err != nil {
			break
		}

		if err = rc.onRecvMsg(msg); err != nil {
			break
		}

	}

	log.Println("rtmp cycle finished, begin clean.err=", err)

	rc.clean()

	log.Println("rtmp clean finished, remote=", rc.tcpConn.Conn.RemoteAddr())
}

func (rc *RtmpConn) clean() {

	log.Println("one publisher begin to quit, stream=", rc.streamName)

	if err := rc.tcpConn.Close(); err != nil {
		log.Println("close socket err=", err)
	}

	if pt.RtmpRoleFlashPublisher == rc.role || pt.RtmpRoleFMLEPublisher == rc.role {
		if nil != rc.source {
			rc.deletePublishStream(rc.streamName)
			log.Println("delete publisher stream=", rc.streamName)
		}
	}

	if pt.RtmpRolePlayer == rc.role {
		if nil != rc.source {
			rc.source.DestroyConsumer(rc.consumer)
		}
	}
}

func (rc *RtmpConn) deletePublishStream(streamName string) {
	gSources.deleteSource(streamName)
}
