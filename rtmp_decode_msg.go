package main

import "fmt"

func (rtmp *RtmpSession) DecodeMsg(chunk *ChunkStruct) (err error) {

	//todo. if recv size > window acknowlegement, send a ack message first.

	//do decode msg.
	switch chunk.msgHeader.msgTypeid {
	case RTMP_MSG_AMF3CommandMessage, RTMP_MSG_AMF0CommandMessage, RTMP_MSG_AMF0DataMessage, RTMP_MSG_AMF3DataMessage:
		err = rtmp.handleAMFCommandAndDataMessage(chunk)
	case RTMP_MSG_UserControlMessage:
	case RTMP_MSG_WindowAcknowledgementSize:
	case RTMP_MSG_SetChunkSize:
	case RTMP_MSG_SetPeerBandwidth:
	case RTMP_MSG_Acknowledgement:
	default:
		err = fmt.Errorf("unknown chunk.msgHeader.msgTypeid=", chunk.msgHeader.msgTypeid)
	}

	if err != nil {
		return
	}

	return
}
