package main

import (
	"UtilsTools/identify_panic"
	"log"
)

func (rtmp *RtmpSession) HandleMsg(chunk *ChunkStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}
	}()

	switch chunk.decodeResultType {
	case DECODE_MSG_TYPE_UNKNOWN:
		log.Println("HandleMsg, unknown msg type.")
	case DECODE_MSG_TYPE_SetChunkSize:
		err = rtmp.HandleMsgSetChunkSize(chunk)
	case DECODE_MSG_TYPE_SetWindowsAcknowlegementSize:
		err = rtmp.HandleMsgSetWindowsAcknowlegementSize(chunk)
	case DECODE_MSG_YTPE_UserControl:
		err = rtmp.HandleMsgUserControl(chunk)
	default:
		log.Println("defalult: HandleMsg unknown msg type.")
	}

	if err != nil {
		return
	}

	return
}

func (rtmp *RtmpSession) HandleMsgSetChunkSize(chunk *ChunkStream) (err error) {
	chunkSize := chunk.decodeResult.(uint32)
	if chunkSize >= RTMP_CHUNKSIZE_MIN && chunkSize <= RTMP_CHUNKSIZE_MAX {
		rtmp.chunkSize = chunkSize
		log.Println("peer set chunk size success. chunk size=", chunkSize)
	} else {
		//ignored
		log.Println("HandleMsgSetChunkSize, chunk size is invalid.", chunkSize)
	}

	return
}

func (rtmp *RtmpSession) HandleMsgSetWindowsAcknowlegementSize(chunk *ChunkStream) (err error) {

	windowAcknowlegementSize := chunk.decodeResult.(uint32)

	if windowAcknowlegementSize > 0 {
		rtmp.ackWindow.ackWindowSize = windowAcknowlegementSize
	} else {
		//ignored.
		log.Println("HandleMsgSetWindowsAcknowlegementSize, ack size is invalied.", windowAcknowlegementSize)
	}

	return
}

func (rtmp *RtmpSession) HandleMsgUserControl(chunk *ChunkStream) (err error) {
	userCtrol := chunk.decodeResult.(UserControlMsg)

	switch userCtrol.eventType {
	case SrcPCUCStreamBegin:
	case SrcPCUCStreamEOF:
	case SrcPCUCStreamDry:
	case SrcPCUCSetBufferLength:
	case SrcPCUCStreamIsRecorded:
	case SrcPCUCPingRequest:
		err = rtmp.ResponsePingMessage(chunk, &userCtrol)
	case SrcPCUCPingResponse:
	default:
		log.Println("HandleMsgUserControl unknown event type.type=", userCtrol.eventType)
	}

	if err != nil {
		return
	}
	return
}
