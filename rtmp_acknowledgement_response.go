package main

import "encoding/binary"

func (rtmp *RtmpSession) AcknowledgementResponse(chunk *ChunkStream) (err error) {

	var msg MessageStream

	//msg payload
	msg.payload = make([]uint8, 4)
	sequenceNum := uint32(rtmp.recvBytesSum)
	binary.BigEndian.PutUint32(msg.payload[:], sequenceNum)

	//msg header
	msg.header.length = 4
	msg.header.typeId = RTMP_MSG_Acknowledgement
	msg.header.streamId = 0
	msg.header.preferCsId = chunk.msg.header.preferCsId

	//begin to send.
	if msg.header.preferCsId < 2 {
		msg.header.preferCsId = RTMP_CID_ProtocolControl
	}

	err = rtmp.SendMsg(&msg)
	if err != nil {
		return
	}

	return
}
