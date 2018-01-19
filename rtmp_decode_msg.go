package main

import (
	"UtilsTools/identify_panic"
	"fmt"
	"log"
)

func (rtmp *RtmpSession) DecodeMsg(chunk *ChunkStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}
	}()

	if (rtmp.ackWindow.ackWindowSize > 0) && (rtmp.recvBytesSum-rtmp.ackWindow.hasAckedSize > uint64(rtmp.ackWindow.ackWindowSize)) {

		err = rtmp.AcknowledgementResponse(chunk)
		if err != nil {
			return
		}

		rtmp.ackWindow.hasAckedSize = rtmp.recvBytesSum
	}

	//do decode msg.
	switch chunk.msg.header.typeId {
	case RTMP_MSG_AMF3CommandMessage, RTMP_MSG_AMF0CommandMessage, RTMP_MSG_AMF0DataMessage, RTMP_MSG_AMF3DataMessage:
		err = rtmp.handleAMFCommandAndDataMessage(chunk)
	case RTMP_MSG_UserControlMessage:
	case RTMP_MSG_WindowAcknowledgementSize:
	case RTMP_MSG_SetChunkSize:
	case RTMP_MSG_SetPeerBandwidth:
	case RTMP_MSG_Acknowledgement:
	default:
		err = fmt.Errorf("unknown chunk.header.typeId=", chunk.msg.header.typeId)
	}

	if err != nil {
		return
	}

	return
}
