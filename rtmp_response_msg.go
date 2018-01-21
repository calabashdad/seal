package main

import "encoding/binary"

func (rtmp *RtmpSession) ResponseConnectApp(chunk *ChunkStream) (err error) {
	var msg MessageStream

	//msg payload
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

	//msg header
	msg.header.length = uint32(len(msg.payload))
	msg.header.typeId = RTMP_MSG_AMF0CommandMessage
	msg.header.streamId = 0
	if chunk.msg.header.preferCsId < 2 {
		msg.header.preferCsId = RTMP_CID_ProtocolControl
	} else {
		msg.header.preferCsId = chunk.msg.header.preferCsId
	}

	err = rtmp.SendMsg(&msg)
	if err != nil {
		return
	}

	return
}

func (rtmp *RtmpSession) ResponsePingMessage(chunk *ChunkStream, userCtrl *UserControlMsg) (err error) {
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
	if chunk.msg.header.preferCsId < 2 {
		msg.header.preferCsId = RTMP_CID_ProtocolControl
	} else {
		msg.header.preferCsId = chunk.msg.header.preferCsId
	}

	err = rtmp.SendMsg(&msg)
	if err != nil {
		return
	}

	return
}
