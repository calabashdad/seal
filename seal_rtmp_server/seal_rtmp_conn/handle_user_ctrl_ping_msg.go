package seal_rtmp_conn

import (
	"encoding/binary"
	"seal/seal_rtmp_server/seal_rtmp_protocol/protocol_stack"
)

func (rtmp *RtmpConn) handleUserCtrlResponsePingMessage(chunkStreamId uint32, userCtrl *UserControlMsg) (err error) {
	var msg MessageStream

	//msg payload
	var offset uint32

	msg.payload = make([]uint8, 2+4) // 2(type) + 4(data)
	binary.BigEndian.PutUint16(msg.payload[offset:offset+2], protocol_stack.SrcPCUCPingResponse)
	offset += 2
	binary.BigEndian.PutUint32(msg.payload[offset:offset+4], userCtrl.eventData)
	offset += 4

	//msg header
	msg.header.length = uint32(len(msg.payload))
	msg.header.typeId = protocol_stack.RTMP_MSG_UserControlMessage
	msg.header.streamId = 0
	if chunkStreamId < 2 {
		msg.header.preferCsId = protocol_stack.RTMP_CID_ProtocolControl
	} else {
		msg.header.preferCsId = chunkStreamId
	}

	err = rtmp.SendMsg(&msg)
	if err != nil {
		return
	}

	return
}
