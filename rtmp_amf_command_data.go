package main

import (
	"UtilsTools/identify_panic"
	"fmt"
	"log"
)

func (rtmp *RtmpSession) handleAMFCommandAndDataMessage(chunk *ChunkStream) (err error) {
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

	if RTMP_AMF0_COMMAND_RESULT == command || RTMP_AMF0_COMMAND_ERROR == command {
		//todo.
	}

	switch command {
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
