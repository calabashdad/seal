package main

import (
	"UtilsTools/identify_panic"
	"encoding/binary"
	"fmt"
	"log"
)

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

func (rtmp *RtmpConn) handleUserControlMessage(msg *MessageStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	log.Println("msg type usercontrol msg.")

	var offset uint32

	if uint32(len(msg.payload))-offset < (2 + 4) {
		err = fmt.Errorf("handleUserControlMessage, 0, length is not enough.")
		return
	}

	var userCtrlMsg UserControlMsg

	userCtrlMsg.eventType = binary.BigEndian.Uint16(msg.payload[offset : offset+2])
	offset += 2

	userCtrlMsg.eventData = binary.BigEndian.Uint32(msg.payload[offset : offset+4])
	offset += 4

	log.Println("user control msg, type=", userCtrlMsg.eventType)

	if SrcPCUCSetBufferLength == userCtrlMsg.eventType {
		if uint32(len(msg.payload))-offset < 4 {
			err = fmt.Errorf("handleUserControlMessage, 1, length is not enough.")
			return
		}

		userCtrlMsg.extraData = binary.BigEndian.Uint32(msg.payload[offset : offset+4])
		offset += 4
	}

	if err != nil {
		return
	}

	switch userCtrlMsg.eventType {
	case SrcPCUCStreamBegin:
	case SrcPCUCStreamEOF:
	case SrcPCUCStreamDry:
	case SrcPCUCSetBufferLength:
	case SrcPCUCStreamIsRecorded:
	case SrcPCUCPingRequest:
		err = rtmp.ResponsePingMessage(msg.header.preferCsId, &userCtrlMsg)
	case SrcPCUCPingResponse:
	default:
		log.Println("HandleMsgUserControl unknown event type.type=", userCtrlMsg.eventType)
	}

	if err != nil {
		return
	}

	return
}

func (rtmp *RtmpConn) ResponsePingMessage(chunkStreamId uint32, userCtrl *UserControlMsg) (err error) {
	var msg MessageStream

	//msg payload
	var offset uint32

	msg.payload = make([]uint8, 2+4) // 2(type) + 4(data)
	binary.BigEndian.PutUint16(msg.payload[offset:offset+2], SrcPCUCPingResponse)
	offset += 2
	binary.BigEndian.PutUint32(msg.payload[offset:offset+4], userCtrl.eventData)
	offset += 4

	//msg header
	msg.header.length = uint32(len(msg.payload))
	msg.header.typeId = RTMP_MSG_UserControlMessage
	msg.header.streamId = 0
	if chunkStreamId < 2 {
		msg.header.preferCsId = RTMP_CID_ProtocolControl
	} else {
		msg.header.preferCsId = chunkStreamId
	}

	err = rtmp.SendMsg(&msg)
	if err != nil {
		return
	}

	return
}
