package seal_rtmp_conn

import (
	"UtilsTools/identify_panic"
	"log"
	"net"
	"sync"
)

//rtmp conn role.
const (
	RTMP_ROLE_UNKNOWN = 0

	RTMP_ROLE_PUBLISH = 1
	RTMP_ROLE_PALY    = 2
)

//key: publish stream without token. value: RtmpConn.
//this map just use to judge if stream has published
var MapPublishingStreams sync.Map

type RtmpConn struct {
	net.Conn
	TimeOut        uint32
	Chunks         map[uint32]*ChunkStream //key csid.
	TransactionIds map[float64]string      //key transaction id. value: request command name
	AckWindow      struct {
		ackWindowSize uint32 //
		hasAckedSize  uint64 //size has acked to peer
	}
	RecvBytesSum   uint64
	ChunkSize      uint32 //default is RTMP_DEFAULT_CHUNK_SIZE. can set by peer.
	Role           uint8  //publish or play. RTMP_ROLE_*
	ObjectEncoding float64
	MetaData       struct {
		marker uint8
		value  interface{}
	}
	StreamInfo struct {
		stream string //withou token.
		token  string
	}

	//key: client remoteAddr. value: chan *MessageStream
	//when role is publish, this is significative
	players sync.Map

	//when role is player, this is significative
	msgChan chan *MessageStream

	//cache
	cacheMsgH264SequenceKeyFrame *MessageStream
	cacheMsgAACSequenceHeader    *MessageStream
	cacheMsgMetaData             *MessageStream
}

func (rtmpSession *RtmpConn) RtmpMsgLoop() (err error) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ", panic at ", identify_panic.IdentifyPanic())
		}
	}()

	for {
		var chunkStreamId uint32
		err, chunkStreamId = rtmpSession.RecvMsg()
		if err != nil {
			break
		}

		err = rtmpSession.HandleMsg(chunkStreamId)
		if err != nil {
			break
		}
	}

	if err != nil {
		return
	}

	return
}
