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

func (rtmp *RtmpSession) ResponseConnectApp(chunk *ChunkStream) (err error) {
	var msg MessageStream

	//msg payload
	if true {
		msg.payload = append(msg.payload, Amf0WriteString(RTMP_AMF0_COMMAND_RESULT)...)
		msg.payload = append(msg.payload, Amf0WriteNumber(1.0)...)

		var objs []Amf0Object

		objs = append(objs, Amf0Object{
			propertyName: "fmsVer",
			value:        "FMS/" + SEAL_VERSION,
			valueType:    RTMP_AMF0_String,
		})

		objs = append(objs, Amf0Object{
			propertyName: "capabilities",
			value:        127.0,
			valueType:    RTMP_AMF0_Number,
		})

		objs = append(objs, Amf0Object{
			propertyName: "mode",
			value:        1.0,
			valueType:    RTMP_AMF0_Number,
		})

		objs = append(objs, Amf0Object{
			propertyName: StatusLevel,
			value:        StatusLevelStatus,
			valueType:    RTMP_AMF0_String,
		})

		objs = append(objs, Amf0Object{
			propertyName: StatusCode,
			value:        StatusCodeConnectSuccess,
			valueType:    RTMP_AMF0_String,
		})

		objs = append(objs, Amf0Object{
			propertyName: StatusDescription,
			value:        "Connection succeeded",
			valueType:    RTMP_AMF0_String,
		})

		objs = append(objs, Amf0Object{
			propertyName: "objectEncoding",
			value:        rtmp.objectEncoding,
			valueType:    RTMP_AMF0_Number,
		})

		msg.payload = append(msg.payload, Amf0WriteObject(objs)...)

		var ecma []Amf0Object

		ecma = append(ecma, Amf0Object{
			propertyName: "version",
			value:        SEAL_VERSION,
			valueType:    RTMP_AMF0_String,
		})

		ecma = append(ecma, Amf0Object{
			propertyName: "seal_license",
			value:        "The MIT License (MIT)",
			valueType:    RTMP_AMF0_String,
		})

		ecma = append(ecma, Amf0Object{
			propertyName: "seal_authors",
			value:        "YangKai",
			valueType:    RTMP_AMF0_String,
		})

		ecma = append(ecma, Amf0Object{
			propertyName: "seal_email",
			value:        "beyondyangkai@gmail.com",
			valueType:    RTMP_AMF0_String,
		})

		ecma = append(ecma, Amf0Object{
			propertyName: "seal_copyright",
			value:        "Copyright (c) 2018 YangKai",
			valueType:    RTMP_AMF0_String,
		})

		ecma = append(ecma, Amf0Object{
			propertyName: "seal_sig",
			value:        "seal",
			valueType:    RTMP_AMF0_String,
		})

		msg.payload = append(msg.payload, Amf0WriteObject(ecma)...)
	}

	//msg header
	msg.header.length = 4
	msg.header.typeId = RTMP_MSG_AMF0CommandMessage
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
