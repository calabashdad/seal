package seal_rtmp_conn

import (
	"UtilsTools/identify_panic"
	"encoding/binary"
	"fmt"
	"log"
	"seal/seal_rtmp_server/seal_rtmp_protocol/protocol_stack"
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

	if protocol_stack.SrcPCUCSetBufferLength == userCtrlMsg.eventType {
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
	case protocol_stack.SrcPCUCStreamBegin:
	case protocol_stack.SrcPCUCStreamEOF:
	case protocol_stack.SrcPCUCStreamDry:
	case protocol_stack.SrcPCUCSetBufferLength:
		err = rtmp.handleUserCtrlSetBufferLength(msg.header.preferCsId, &userCtrlMsg)
	case protocol_stack.SrcPCUCStreamIsRecorded:
	case protocol_stack.SrcPCUCPingRequest:
		err = rtmp.handleUserCtrlResponsePingMessage(msg.header.preferCsId, &userCtrlMsg)
	case protocol_stack.SrcPCUCPingResponse:
	default:
		log.Println("HandleMsgUserControl unknown event type.type=", userCtrlMsg.eventType)
	}

	if err != nil {
		return
	}

	return
}
