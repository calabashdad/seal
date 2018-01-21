package main

import (
	"UtilsTools/identify_panic"
	"encoding/binary"
	"fmt"
	"log"
)

const (
	DECODE_MSG_TYPE_UNKNOWN                      = 0
	DECODE_MSG_TYPE_Amf0CommandConnectPkg        = 1
	DECODE_MSG_TYPE_SetChunkSize                 = 2
	DECODE_MSG_TYPE_SetWindowsAcknowlegementSize = 3
	DECODE_MSG_YTPE_UserControl                  = 4
)

func (rtmp *RtmpSession) DecodeMsg(chunk *ChunkStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}
	}()

	//reset the decode type. incase the multiplex channel.
	chunk.decodeResultType = DECODE_MSG_TYPE_UNKNOWN

	switch chunk.msg.header.typeId {
	case RTMP_MSG_AMF3CommandMessage, RTMP_MSG_AMF0CommandMessage, RTMP_MSG_AMF0DataMessage, RTMP_MSG_AMF3DataMessage:
		err = rtmp.decodeAMFCommandAndDataMessage(chunk)
	case RTMP_MSG_UserControlMessage:
		err = rtmp.decodeUserControlMessage(chunk)
	case RTMP_MSG_WindowAcknowledgementSize:
		err = rtmp.decodeSetWindowAcknowledgementSize(chunk)
	case RTMP_MSG_SetChunkSize:
		err = rtmp.decodeSetChunkSize(chunk)
	case RTMP_MSG_SetPeerBandwidth:
		err = rtmp.decodeSetPeerBandWidth(chunk)
	case RTMP_MSG_Acknowledgement:
		err = rtmp.decodeAcknowlegement(chunk)
	default:
		err = fmt.Errorf("unknown chunk.header.typeId=", chunk.msg.header.typeId)
	}

	if err != nil {
		return
	}

	return
}

func (rtmp *RtmpSession) decodeAMFCommandAndDataMessage(chunk *ChunkStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}
	}()

	var offset uint32

	// skip 1bytes to decode the amf3 command.
	if RTMP_MSG_AMF3CommandMessage == chunk.msg.header.typeId && chunk.msg.header.length >= 1 {
		offset += 1
	}

	//read the command name.
	err, command := Amf0ReadString(chunk.msg.payload, &offset)
	if err != nil {
		return
	}

	if 0 == len(command) {
		err = fmt.Errorf("Amf0ReadString failed, command is nil.")
		return
	}

	switch command {
	case RTMP_AMF0_COMMAND_RESULT, RTMP_AMF0_COMMAND_ERROR:
		err = rtmp.handleAMF0CommandResultError(chunk)
	case RTMP_AMF0_COMMAND_CONNECT:
		err = rtmp.handleAMF0CommandConnect(chunk)
	case RTMP_AMF0_COMMAND_CREATE_STREAM:
		//todo.
	case RTMP_AMF0_COMMAND_PLAY:
		//todo.
	case RTMP_AMF0_COMMAND_PAUSE:
		//todo.
	case RTMP_AMF0_COMMAND_RELEASE_STREAM:
		//todo.
	case RTMP_AMF0_COMMAND_FC_PUBLISH:
		//todo.
	case RTMP_AMF0_COMMAND_PUBLISH:
		//todo.
	case RTMP_AMF0_COMMAND_UNPUBLISH:
		//todo.
	case RTMP_AMF0_COMMAND_KEEPLIVE:
		//todo.
	case RTMP_AMF0_DATA_SET_DATAFRAME, RTMP_AMF0_DATA_ON_METADATA:
		//todo.
	case RTMP_AMF0_DATA_ON_CUSTOMDATA:
		//todo.
	case SRS_BW_CHECK_FINISHED, SRS_BW_CHECK_PLAYING, SRS_BW_CHECK_PUBLISHING,
		SRS_BW_CHECK_STARTING_PLAY, SRS_BW_CHECK_STARTING_PUBLISH, SRS_BW_CHECK_START_PLAY,
		SRS_BW_CHECK_START_PUBLISH, SRS_BW_CHECK_STOPPED_PLAY, SRS_BW_CHECK_STOP_PLAY,
		SRS_BW_CHECK_STOP_PUBLISH, SRS_BW_CHECK_STOPPED_PUBLISH, SRS_BW_CHECK_FINAL:
		//todo.
	case RTMP_AMF0_COMMAND_CLOSE_STREAM:
		//todo.
	default:
		if RTMP_MSG_AMF0CommandMessage == chunk.msg.header.typeId || RTMP_MSG_AMF3CommandMessage == chunk.msg.header.typeId {
			//todo.
		} else {
			err = fmt.Errorf("unknown command type, command=", command)
		}
	}

	if err != nil {
		return
	}

	return
}

/**
* 5.4. User Control Message (4)
*
* for the EventData is 4bytes.
* Stream Begin(=0)              4-bytes stream ID
* Stream EOF(=1)                4-bytes stream ID
* StreamDry(=2)                 4-bytes stream ID
* SetBufferLength(=3)           8-bytes 4bytes stream ID, 4bytes buffer length.
* StreamIsRecorded(=4)          4-bytes stream ID
* PingRequest(=6)               4-bytes timestamp local server time
* PingResponse(=7)              4-bytes timestamp received ping request.
*
* 3.7. User Control message
* +------------------------------+-------------------------
* | Event Type ( 2- bytes ) | Event Data
* +------------------------------+-------------------------
* Figure 5 Pay load for the ‘User Control Message’.
 */
type UserControlMsg struct {
	eventType uint16
	eventData uint32
	/**
	 * 4bytes if event_type is SetBufferLength; otherwise 0.
	 */
	extraData uint32
}

func (rtmp *RtmpSession) decodeUserControlMessage(chunk *ChunkStream) (err error) {
	var offset uint32

	if uint32(len(chunk.msg.payload))-offset < (2 + 4) {
		err = fmt.Errorf("decodeUserControlMessage, 0, length is not enough.")
		return
	}

	var userCtrlMsg UserControlMsg

	userCtrlMsg.eventType = binary.BigEndian.Uint16(chunk.msg.payload[offset : offset+2])
	offset += 2

	userCtrlMsg.eventData = binary.BigEndian.Uint32(chunk.msg.payload[offset : offset+4])
	offset += 4

	if SrcPCUCSetBufferLength == userCtrlMsg.eventType {
		if uint32(len(chunk.msg.payload))-offset < 4 {
			err = fmt.Errorf("decodeUserControlMessage, 1, length is not enough.")
			return
		}

		userCtrlMsg.extraData = binary.BigEndian.Uint32(chunk.msg.payload[offset : offset+4])
		offset += 4
	}

	if err != nil {
		return
	}

	chunk.decodeResult = userCtrlMsg
	chunk.decodeResultType = DECODE_MSG_YTPE_UserControl

	return
}

func (rtmp *RtmpSession) decodeSetWindowAcknowledgementSize(chunk *ChunkStream) (err error) {
	if len(chunk.msg.payload) >= 4 {
		chunk.decodeResult = binary.BigEndian.Uint32(chunk.msg.payload[0:4])
		chunk.decodeResultType = DECODE_MSG_TYPE_SetWindowsAcknowlegementSize
	} else {
		err = fmt.Errorf("decodeSetWindowAcknowledgementSize payload len < 4", len(chunk.msg.payload))
		return
	}
	return
}

func (rtmp *RtmpSession) decodeSetChunkSize(chunk *ChunkStream) (err error) {

	if len(chunk.msg.payload) >= 4 {
		chunk.decodeResult = binary.BigEndian.Uint32(chunk.msg.payload[0:4])
		chunk.decodeResultType = DECODE_MSG_TYPE_SetChunkSize
	} else {
		err = fmt.Errorf("decodeSetChunkSize payload length < 4", len(chunk.msg.payload))
		return
	}

	return
}

func (rtmp *RtmpSession) decodeSetPeerBandWidth(chunk *ChunkStream) (err error) {
	return
}

func (rtmp *RtmpSession) decodeAcknowlegement(chunk *ChunkStream) (err error) {
	return
}
