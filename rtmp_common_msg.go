package main

import (
	"encoding/binary"
)

func (rtmp *RtmpSession) CommonMsgSetWindowAcknowledgementSize(chunk *ChunkStream, WindowAcknowledgementSize uint32) (err error) {

	var msg MessageStream

	//msg payload
	msg.payload = make([]uint8, 4)
	binary.BigEndian.PutUint32(msg.payload[:], WindowAcknowledgementSize)

	//msg header
	msg.header.length = 4
	msg.header.typeId = RTMP_MSG_WindowAcknowledgementSize
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

func (rtmp *RtmpSession) CommonMsgResponseWindowAcknowledgement(chunk *ChunkStream, ackedSize uint32) (err error) {

	var msg MessageStream

	//msg payload
	msg.payload = make([]uint8, 4)
	binary.BigEndian.PutUint32(msg.payload[:], ackedSize)

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

func (rtmp *RtmpSession) CommonMsgSetPeerBandwidth(chunk *ChunkStream, bandWidthValue uint32, limitType uint8) (err error) {

	var msg MessageStream

	//msg payload
	msg.payload = make([]uint8, 5)
	binary.BigEndian.PutUint32(msg.payload[:4], bandWidthValue)
	msg.payload[4] = limitType

	//msg header
	msg.header.length = 4
	msg.header.typeId = RTMP_MSG_SetPeerBandwidth
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

func (rtmp *RtmpSession) CommonMsgSetChunkSize(chunkSize uint32) (err error) {
	var msg MessageStream

	//msg payload
	msg.payload = make([]uint8, 4)
	binary.BigEndian.PutUint32(msg.payload[:], chunkSize)

	//msg header
	msg.header.length = 4
	msg.header.typeId = RTMP_MSG_SetChunkSize
	msg.header.streamId = 0
	msg.header.preferCsId = RTMP_CID_ProtocolControl

	err = rtmp.SendMsg(&msg)
	if err != nil {
		return
	}

	return
}
