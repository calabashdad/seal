package seal_rtmp_conn

import (
	"UtilsTools/identify_panic"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"seal/seal_rtmp_server/seal_rtmp_protocol/handshake"
	"seal/seal_rtmp_server/seal_rtmp_protocol/protocol_stack"
	"sync"
)

//rtmp conn role.
const (
	RTMP_ROLE_UNKNOWN = 0

	RTMP_ROLE_PUBLISH = 1
	RTMP_ROLE_PALY    = 2
)

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
	Role           uint8  //publish or play.
	ObjectEncoding float64
	MetaData       struct {
		marker uint8
		value  interface{}
	}
	StreamInfo struct {
		stream string //withou token.
		token  string
	}
}

func (rtmpSession *RtmpConn) HandleRtmpSession() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}

		rtmpSession.Conn.Close()
		MapPublishingStreams.Delete(rtmpSession.StreamInfo.stream)

		log.Println("One RtmpConn finished.remote=", rtmpSession.Conn.RemoteAddr())
	}()

	log.Println("One RtmpConn come in. remote=", rtmpSession.Conn.RemoteAddr())

	err := rtmpSession.HandShake()
	if err != nil {
		log.Println("rtmp handshake failed, err=", err)
		return
	}

	log.Println("rtmp handshake success.remote=", rtmpSession.Conn.RemoteAddr())

	err = rtmpSession.RtmpMsgLoop()

	log.Println("rtmp msg loop quit.err=", err, ",remote=", rtmpSession.Conn.RemoteAddr())
}

func (rtmp *RtmpConn) HandShake() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}
	}()

	var handshakeData [6146]uint8 // c0(1) + c1(1536) + c2(1536) + s0(1) + s1(1536) + s2(1536)

	c0 := handshakeData[:1]
	c1 := handshakeData[1:1537]
	c2 := handshakeData[1537:3073]

	s0 := handshakeData[3073:3074]
	s1 := handshakeData[3074:4610]
	s2 := handshakeData[4610:6146]

	c0c1 := handshakeData[0:1537]
	s0s1s2 := handshakeData[3073:6146]

	//recv c0c1
	err = rtmp.ExpectBytes(1537, c0c1)
	if err != nil {
		return
	}

	//parse c0
	if c0[0] != 3 {
		err = fmt.Errorf("client c0 is not 3.")
		return
	}

	//use complex handshake, if complex handshake failed, try use simple handshake
	//parse c1
	clientVer := binary.BigEndian.Uint32(c1[4:8])
	if 0 != clientVer {
		if !handshake.ComplexHandShake(c1, s0, s1, s2) {
			err = fmt.Errorf("0 != clientVer, complex handshake failed.")
			return
		}
	} else {
		//use simple handshake
		log.Println("0 == clientVer, client use simple handshake.")
		s0[0] = 3
		copy(s1, c2)
		copy(s2, c1)
	}

	//send s0s1s2
	err = rtmp.SendBytes(s0s1s2)
	if err != nil {
		return
	}

	//recv c2
	err = rtmp.ExpectBytes(uint32(len(c2)), c2)
	if err != nil {
		return
	}

	//c2 do not need verify.

	return
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

		err = rtmpSession.HanleMsg(chunkStreamId)
		if err != nil {
			break
		}
	}

	if err != nil {
		return
	}

	return
}

func (rtmp *RtmpConn) HanleMsg(chunkStreamId uint32) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}
	}()

	chunk := rtmp.Chunks[chunkStreamId]
	if nil == chunk {
		err = fmt.Errorf("HanleMsg:can not find the chunk strema id in chuns.")
		return
	}

	switch chunk.msg.header.typeId {
	case protocol_stack.RTMP_MSG_AMF3CommandMessage, protocol_stack.RTMP_MSG_AMF0CommandMessage,
		protocol_stack.RTMP_MSG_AMF0DataMessage, protocol_stack.RTMP_MSG_AMF3DataMessage:
		err = rtmp.handleAMFCommandAndDataMessage(&chunk.msg)
	case protocol_stack.RTMP_MSG_UserControlMessage:
		err = rtmp.handleUserControlMessage(&chunk.msg)
	case protocol_stack.RTMP_MSG_WindowAcknowledgementSize:
		err = rtmp.handleSetWindowAcknowledgementSize(&chunk.msg)
	case protocol_stack.RTMP_MSG_SetChunkSize:
		err = rtmp.handleSetChunkSize(&chunk.msg)
	case protocol_stack.RTMP_MSG_SetPeerBandwidth:
		err = rtmp.handleSetPeerBandWidth(&chunk.msg)
	case protocol_stack.RTMP_MSG_Acknowledgement:
		err = rtmp.handleAcknowlegement(&chunk.msg)
	case protocol_stack.RTMP_MSG_AbortMessage:
		err = rtmp.handleAbortMsg(&chunk.msg)
	case protocol_stack.RTMP_MSG_EdgeAndOriginServerCommand:
		err = rtmp.handleEdgeAndOriginServerCommand(&chunk.msg)
	case protocol_stack.RTMP_MSG_AMF3SharedObject:
		//todo
	case protocol_stack.RTMP_MSG_AMF0SharedObject:
		//todo
	case protocol_stack.RTMP_MSG_AudioMessage:
		err = rtmp.handleMsgAudio(&chunk.msg)
	case protocol_stack.RTMP_MSG_VideoMessage:
		err = rtmp.handleMsgVideo(&chunk.msg)
	case protocol_stack.RTMP_MSG_AggregateMessage:
		//todo.
	default:
		err = fmt.Errorf("HanleMsg: unknown msg type. typeid=", chunk.msg.header.typeId)
	}

	if err != nil {
		return
	}

	return
}
