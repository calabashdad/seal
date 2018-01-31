package co

import (
	"UtilsTools/identify_panic"
	"fmt"
	"log"
	"seal/rtmp/pt"
)

func (rc *RtmpConn) onRecvMsg(csid uint32) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	err = rc.EstimateNeedSendAcknowlegement()
	if err != nil {
		return
	}

	chunk := rc.ChunkStreams[csid]
	if nil == chunk {
		err = fmt.Errorf("on recv msg can not find the csid.")
		return
	}

	switch chunk.Msg.Header.MessageType {
	case pt.RTMP_MSG_AMF3CommandMessage, pt.RTMP_MSG_AMF0CommandMessage,
		pt.RTMP_MSG_AMF0DataMessage, pt.RTMP_MSG_AMF3DataMessage:
		err = rc.msgAmf(&chunk.Msg)
	case pt.RTMP_MSG_UserControlMessage:
		err = rc.msgUserCtrl(&chunk.Msg)
	case pt.RTMP_MSG_WindowAcknowledgementSize:
		err = rc.msgSetAck(&chunk.Msg)
	case pt.RTMP_MSG_SetChunkSize:
		err = rc.msgSetChunk(&chunk.Msg)
	case pt.RTMP_MSG_SetPeerBandwidth:
		err = rc.msgSetBand(&chunk.Msg)
	case pt.RTMP_MSG_Acknowledgement:
		err = rc.msgAck(&chunk.Msg)
	case pt.RTMP_MSG_AbortMessage:
		err = rc.msgAbort(&chunk.Msg)
	case pt.RTMP_MSG_EdgeAndOriginServerCommand:
		//todo
	case pt.RTMP_MSG_AMF3SharedObject:
		//todo
	case pt.RTMP_MSG_AMF0SharedObject:
		//todo
	case pt.RTMP_MSG_AudioMessage:
		err = rc.msgAudio(&chunk.Msg)
	case pt.RTMP_MSG_VideoMessage:
		err = rc.msgVideo(&chunk.Msg)
	case pt.RTMP_MSG_AggregateMessage:
		err = rc.msgAggregate(&chunk.Msg)
	default:
		log.Println("on recv msg unknown msg typeid=", chunk.Msg.Header.MessageType)
	}

	if err != nil {
		return
	}

	return
}
