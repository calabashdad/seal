package co

import (
	"log"
	"seal/rtmp/pt"

	"github.com/calabashdad/utiltools"
)

func (rc *RtmpConn) onRecvMsg(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	if nil == msg {
		return
	}

	if err = rc.estimateNeedSendAcknowlegement(); err != nil {
		return
	}

	switch msg.Header.MessageType {
	case pt.RtmpMsgAmf3CommandMessage, pt.RtmpMsgAmf0CommandMessage,
		pt.RtmpMsgAmf0DataMessage, pt.RtmpMsgAmf3DataMessage:
		err = rc.msgAmf(msg)
	case pt.RtmpMsgUserControlMessage:
		err = rc.msgUserCtrl(msg)
	case pt.RtmpMsgWindowAcknowledgementSize:
		err = rc.msgSetAck(msg)
	case pt.RtmpMsgSetChunkSize:
		err = rc.msgSetChunk(msg)
	case pt.RtmpMsgSetPeerBandwidth:
		err = rc.msgSetBand(msg)
	case pt.RtmpMsgAcknowledgement:
		err = rc.msgAck(msg)
	case pt.RtmpMsgAbortMessage:
		err = rc.msgAbort(msg)
	case pt.RtmpMsgEdgeAndOriginServerCommand:
	case pt.RtmpMsgAmf3SharedObject:
	case pt.RtmpMsgAmf0SharedObject:
	case pt.RtmpMsgAudioMessage:
		err = rc.msgAudio(msg)
	case pt.RtmpMsgVideoMessage:
		err = rc.msgVideo(msg)
	case pt.RtmpMsgAggregateMessage:
		err = rc.msgAggregate(msg)
	default:
		log.Println("on recv msg unknown msg typeid=", msg.Header.MessageType)
	}

	if err != nil {
		return
	}

	return
}
