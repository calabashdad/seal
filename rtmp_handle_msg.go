package main

import (
	"UtilsTools/identify_panic"
	"fmt"
	"log"
)

func (rtmp *RtmpConn) HanleMsg(chunkStreamId uint32) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}
	}()

	chunk := rtmp.chunks[chunkStreamId]
	if nil == chunk {
		err = fmt.Errorf("HanleMsg:can not find the chunk strema id in chuns.")
		return
	}

	switch chunk.msg.header.typeId {
	case RTMP_MSG_AMF3CommandMessage, RTMP_MSG_AMF0CommandMessage, RTMP_MSG_AMF0DataMessage, RTMP_MSG_AMF3DataMessage:
		err = rtmp.handleAMFCommandAndDataMessage(&chunk.msg)
	case RTMP_MSG_UserControlMessage:
		err = rtmp.handleUserControlMessage(&chunk.msg)
	case RTMP_MSG_WindowAcknowledgementSize:
		err = rtmp.handleSetWindowAcknowledgementSize(&chunk.msg)
	case RTMP_MSG_SetChunkSize:
		err = rtmp.handleSetChunkSize(&chunk.msg)
	case RTMP_MSG_SetPeerBandwidth:
		err = rtmp.handleSetPeerBandWidth(&chunk.msg)
	case RTMP_MSG_Acknowledgement:
		err = rtmp.handleAcknowlegement(&chunk.msg)
	case RTMP_MSG_AbortMessage:
		err = rtmp.handleAbortMsg(&chunk.msg)
	case RTMP_MSG_EdgeAndOriginServerCommand:
		err = rtmp.handleEdgeAndOriginServerCommand(&chunk.msg)
	default:
		err = fmt.Errorf("HanleMsg: unknown msg type. ", chunk.msg.header.typeId)
	}

	if err != nil {
		return
	}

	return
}
