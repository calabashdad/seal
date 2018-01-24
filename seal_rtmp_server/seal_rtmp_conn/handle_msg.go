package seal_rtmp_conn

import (
	"UtilsTools/identify_panic"
	"fmt"
	"log"
	"seal/seal_rtmp_server/seal_rtmp_protocol/protocol_stack"
)

func (rtmp *RtmpConn) HandleMsg(chunkStreamId uint32) (err error) {
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
		log.Println("HanleMsg: unknown msg type. typeid=", chunk.msg.header.typeId)
	}

	if err != nil {
		return
	}

	return
}
