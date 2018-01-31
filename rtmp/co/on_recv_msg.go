package co

import (
	"UtilsTools/identify_panic"
	"fmt"
	"log"
	"seal/rtmp/pt"
)

func (rc *RtmpConn) OnRecvMsg(csid uint32) (err error) {
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
		err = rc.MsgAmf(&chunk.Msg)
	case pt.RTMP_MSG_UserControlMessage:
		err = rc.MsgUserCtrl(&chunk.Msg)
	case pt.RTMP_MSG_WindowAcknowledgementSize:
		err = rc.MsgSetAck(&chunk.Msg)
	case pt.RTMP_MSG_SetChunkSize:
		err = rc.MsgSetChunk(&chunk.Msg)
	case pt.RTMP_MSG_SetPeerBandwidth:
		err = rc.MsgSetBand(&chunk.Msg)
	case pt.RTMP_MSG_Acknowledgement:
		err = rc.MsgAck(&chunk.Msg)
	case pt.RTMP_MSG_AbortMessage:
		err = rc.MsgAbort(&chunk.Msg)
	case pt.RTMP_MSG_EdgeAndOriginServerCommand:
		//todo
	case pt.RTMP_MSG_AMF3SharedObject:
		//todo
	case pt.RTMP_MSG_AMF0SharedObject:
		//todo
	case pt.RTMP_MSG_AudioMessage:
		err = rc.MsgAudio(&chunk.Msg)
	case pt.RTMP_MSG_VideoMessage:
		err = rc.MsgVideo(&chunk.Msg)
	case pt.RTMP_MSG_AggregateMessage:
		err = rc.MsgAggregate(&chunk.Msg)
	default:
		log.Println("on recv msg unknown msg typeid=", chunk.Msg.Header.MessageType)
	}

	if err != nil {
		return
	}

	return
}
