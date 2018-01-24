package seal_rtmp_conn

import (
	"UtilsTools/identify_panic"
	"encoding/binary"
	"fmt"
	"log"
	"seal/seal_rtmp_server/seal_rtmp_protocol/amf_serial"
	"seal/seal_rtmp_server/seal_rtmp_protocol/protocol_stack"
	"strconv"
	"strings"
)

func (rtmp *RtmpConn) handleAMF0CmdConnect(msg *MessageStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	var connectPkg Amf0CommandConnectPkg

	var offset uint32

	err, connectPkg.command = amf_serial.Amf0ReadString(msg.payload, &offset)
	if err != nil {
		return
	}

	if connectPkg.command != protocol_stack.RTMP_AMF0_COMMAND_CONNECT {
		err = fmt.Errorf("handleAMF0CmdConnect command is error. command=", connectPkg.command)
		return
	}

	err, connectPkg.transactionId = amf_serial.Amf0ReadNumber(msg.payload, &offset)
	if err != nil {
		return
	}

	//this method is not strict for float type. just a warn.
	if 1 != connectPkg.transactionId {
		log.Println("warn:handleAMF0CmdConnect: transactionId is not 1. transactionId=", connectPkg.transactionId)
	}

	err, connectPkg.commandObjects = amf_serial.Amf0ReadObject(msg.payload, &offset)
	if err != nil {
		return
	}

	if offset < uint32(len(msg.payload)) {
		var v interface{}
		var marker uint8
		err, v, marker = amf_serial.Amf0Discovery(msg.payload, &offset)
		if err != nil {
			return
		}

		if protocol_stack.RTMP_AMF0_Object == marker {
			connectPkg.amfOptional = v
		}
	}

	err = rtmp.ParseConnectPkg(&connectPkg)
	if err != nil {
		log.Println("parse connect pkg error.", err)
		return
	}

	err = rtmp.CommonMsgSetWindowAcknowledgementSize(msg.header.preferCsId, 2500000)
	if err != nil {
		return
	}

	err = rtmp.CommonMsgSetPeerBandwidth(msg.header.preferCsId, 2500000, 2)
	if err != nil {
		return
	}

	err = rtmp.ResponseConnectApp(msg.header.preferCsId)
	if err != nil {
		return
	}

	log.Println("handle amf0 cmd connect success.")

	return
}

type Amf0CommandConnectPkg struct {
	command        string
	transactionId  float64
	commandObjects []amf_serial.Amf0Object
	amfOptional    interface{}
}

func (pkg *Amf0CommandConnectPkg) GetProperty(key string) (value interface{}) {

	for _, v := range pkg.commandObjects {
		if v.PropertyName == key {
			return v.Value
		}
	}

	return
}

func (rtmp *RtmpConn) ResponseConnectApp(chunkStreamId uint32) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	var msg MessageStream

	//msg payload
	msg.payload = append(msg.payload, amf_serial.Amf0WriteString(protocol_stack.RTMP_AMF0_COMMAND_RESULT)...)
	msg.payload = append(msg.payload, amf_serial.Amf0WriteNumber(1.0)...) //transaction id

	var objs []amf_serial.Amf0Object

	objs = append(objs, amf_serial.Amf0Object{
		PropertyName: "fmsVer",
		Value:        "FMS/" + protocol_stack.FMS_VERSION,
		ValueType:    protocol_stack.RTMP_AMF0_String,
	})

	objs = append(objs, amf_serial.Amf0Object{
		PropertyName: "capabilities",
		Value:        127.0,
		ValueType:    protocol_stack.RTMP_AMF0_Number,
	})

	objs = append(objs, amf_serial.Amf0Object{
		PropertyName: "mode",
		Value:        1.0,
		ValueType:    protocol_stack.RTMP_AMF0_Number,
	})

	objs = append(objs, amf_serial.Amf0Object{
		PropertyName: protocol_stack.StatusLevel,
		Value:        protocol_stack.StatusLevelStatus,
		ValueType:    protocol_stack.RTMP_AMF0_String,
	})

	objs = append(objs, amf_serial.Amf0Object{
		PropertyName: protocol_stack.StatusCode,
		Value:        protocol_stack.StatusCodeConnectSuccess,
		ValueType:    protocol_stack.RTMP_AMF0_String,
	})

	objs = append(objs, amf_serial.Amf0Object{
		PropertyName: protocol_stack.StatusDescription,
		Value:        "Connection succeeded",
		ValueType:    protocol_stack.RTMP_AMF0_String,
	})

	objs = append(objs, amf_serial.Amf0Object{
		PropertyName: "objectEncoding",
		Value:        rtmp.ObjectEncoding,
		ValueType:    protocol_stack.RTMP_AMF0_Number,
	})

	msg.payload = append(msg.payload, amf_serial.Amf0WriteObject(objs)...)

	var ecma []amf_serial.Amf0Object

	ecma = append(ecma, amf_serial.Amf0Object{
		PropertyName: "version",
		Value:        protocol_stack.FMS_VERSION,
		ValueType:    protocol_stack.RTMP_AMF0_String,
	})

	ecma = append(ecma, amf_serial.Amf0Object{
		PropertyName: "seal_license",
		Value:        "The MIT License (MIT)",
		ValueType:    protocol_stack.RTMP_AMF0_String,
	})

	ecma = append(ecma, amf_serial.Amf0Object{
		PropertyName: "seal_authors",
		Value:        "YangKai",
		ValueType:    protocol_stack.RTMP_AMF0_String,
	})

	ecma = append(ecma, amf_serial.Amf0Object{
		PropertyName: "seal_email",
		Value:        "beyondyangkai@gmail.com",
		ValueType:    protocol_stack.RTMP_AMF0_String,
	})

	ecma = append(ecma, amf_serial.Amf0Object{
		PropertyName: "seal_copyright",
		Value:        "Copyright (c) 2018 YangKai",
		ValueType:    protocol_stack.RTMP_AMF0_String,
	})

	ecma = append(ecma, amf_serial.Amf0Object{
		PropertyName: "seal_sig",
		Value:        "seal",
		ValueType:    protocol_stack.RTMP_AMF0_String,
	})

	msg.payload = append(msg.payload, amf_serial.Amf0WriteObject(ecma)...)

	//msg header
	msg.header.length = uint32(len(msg.payload))
	msg.header.typeId = protocol_stack.RTMP_MSG_AMF0CommandMessage
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

func (rtmp *RtmpConn) ParseConnectPkg(pkg *Amf0CommandConnectPkg) (err error) {
	tcUrlValue := pkg.GetProperty("tcUrl")
	if nil == tcUrlValue {
		err = fmt.Errorf("tcUrl is nil.")
		return
	}

	objectEncodingValue := pkg.GetProperty("objectEncoding")
	if objectEncodingValue != nil {
		rtmp.ObjectEncoding = objectEncodingValue.(float64)
	}

	var rtmpUrlData RtmpUrlData
	err = rtmpUrlData.ParseUrl(tcUrlValue.(string))
	if err != nil {
		return
	}

	err = rtmpUrlData.Discover()
	if err != nil {
		return
	}

	return
}

type RtmpUrlData struct {
	schema string
	host   string
	port   uint16
	app    string
	stream string
	token  string
}

//format: rtmp://127.0.0.1:1935/live/test?token=abc123
func (urlData *RtmpUrlData) ParseUrl(url string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}
	}()

	//url is not contain the stream with token.
	var urlTmp string

	urlTmp = strings.Replace(url, "://", " ", 1)
	urlTmp = strings.Replace(urlTmp, ":", " ", 1)
	urlTmp = strings.Replace(urlTmp, "/", " ", 1)

	urlSplit := strings.Split(urlTmp, " ")

	if 4 == len(urlSplit) {
		urlData.schema = urlSplit[0]
		urlData.host = urlSplit[1]
		port, ok := strconv.Atoi(urlSplit[2])
		if nil == ok {
			//the port is not default
			if port > 0 && port < 65536 {
				urlData.port = uint16(port)
			} else {
				err = fmt.Errorf("tcUrl port format is error, port=", port)
				return
			}

		} else {
			err = fmt.Errorf("tcurl format error when convert port format, err=", ok)
			return
		}
		urlData.app = urlSplit[3]
	} else {
		err = fmt.Errorf("tcUrl format is error. tcUrl=", url)
		return
	}

	return
}

func (urlData *RtmpUrlData) Discover() (err error) {
	if 0 == len(urlData.schema) ||
		0 == len(urlData.host) ||
		0 == len(urlData.app) {
		err = fmt.Errorf("discover url data failed. url data=", urlData)
		return
	}
	return
}

func (rtmp *RtmpConn) CommonMsgSetWindowAcknowledgementSize(chunkStreamId uint32, WindowAcknowledgementSize uint32) (err error) {

	var msg MessageStream

	//msg payload
	msg.payload = make([]uint8, 4)
	binary.BigEndian.PutUint32(msg.payload[:], WindowAcknowledgementSize)

	//msg header
	msg.header.length = 4
	msg.header.typeId = protocol_stack.RTMP_MSG_WindowAcknowledgementSize
	msg.header.streamId = 0
	msg.header.preferCsId = chunkStreamId

	//begin to send.
	if msg.header.preferCsId < 2 {
		msg.header.preferCsId = protocol_stack.RTMP_CID_ProtocolControl
	}

	err = rtmp.SendMsg(&msg)
	if err != nil {
		return
	}

	return
}

func (rtmp *RtmpConn) CommonMsgSetPeerBandwidth(chunkStreamId uint32, bandWidthValue uint32, limitType uint8) (err error) {

	var msg MessageStream

	//msg payload
	msg.payload = make([]uint8, 5)
	binary.BigEndian.PutUint32(msg.payload[:4], bandWidthValue)
	msg.payload[4] = limitType

	//msg header
	msg.header.length = 4
	msg.header.typeId = protocol_stack.RTMP_MSG_SetPeerBandwidth
	msg.header.streamId = 0
	msg.header.preferCsId = chunkStreamId

	//begin to send.
	if msg.header.preferCsId < 2 {
		msg.header.preferCsId = protocol_stack.RTMP_CID_ProtocolControl
	}

	err = rtmp.SendMsg(&msg)
	if err != nil {
		return
	}

	return
}

func (rtmp *RtmpConn) CommonMsgSetChunkSize(chunkSize uint32) (err error) {
	var msg MessageStream

	//msg payload
	msg.payload = make([]uint8, 4)
	binary.BigEndian.PutUint32(msg.payload[:], chunkSize)

	//msg header
	msg.header.length = 4
	msg.header.typeId = protocol_stack.RTMP_MSG_SetChunkSize
	msg.header.streamId = 0
	msg.header.preferCsId = protocol_stack.RTMP_CID_ProtocolControl

	err = rtmp.SendMsg(&msg)
	if err != nil {
		return
	}

	return
}
