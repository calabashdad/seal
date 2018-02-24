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

	err = rc.EstimateNeedSendAcknowlegement()
	if err != nil {
		return
	}

	if nil == msg {

		return
	}

	switch msg.Header.MessageType {
	case pt.RTMP_MSG_AMF3CommandMessage, pt.RTMP_MSG_AMF0CommandMessage,
		pt.RTMP_MSG_AMF0DataMessage, pt.RTMP_MSG_AMF3DataMessage:
		err = rc.msgAmf(msg)
	case pt.RTMP_MSG_UserControlMessage:
		err = rc.msgUserCtrl(msg)
	case pt.RTMP_MSG_WindowAcknowledgementSize:
		err = rc.msgSetAck(msg)
	case pt.RTMP_MSG_SetChunkSize:
		err = rc.msgSetChunk(msg)
	case pt.RTMP_MSG_SetPeerBandwidth:
		err = rc.msgSetBand(msg)
	case pt.RTMP_MSG_Acknowledgement:
		err = rc.msgAck(msg)
	case pt.RTMP_MSG_AbortMessage:
		err = rc.msgAbort(msg)
	case pt.RTMP_MSG_EdgeAndOriginServerCommand:
	case pt.RTMP_MSG_AMF3SharedObject:
	case pt.RTMP_MSG_AMF0SharedObject:
	case pt.RTMP_MSG_AudioMessage:
		err = rc.msgAudio(msg)
	case pt.RTMP_MSG_VideoMessage:
		err = rc.msgVideo(msg)
	case pt.RTMP_MSG_AggregateMessage:
		err = rc.msgAggregate(msg)
	default:
		log.Println("on recv msg unknown msg typeid=", msg.Header.MessageType)
	}

	if err != nil {
		return
	}

	return
}
