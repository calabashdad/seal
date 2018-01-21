package main

import (
	"UtilsTools/identify_panic"
	"fmt"
	"log"
)

type Amf0CommandConnectPkg struct {
	command        string
	transactionId  float64
	commandObjects []Amf0Object
	amfOptional    interface{}
}

func (rtmp *RtmpConn) handleAMF0CommandConnect(msg *MessageStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	var connectPkg Amf0CommandConnectPkg

	var offset uint32

	err, connectPkg.command = Amf0ReadString(msg.payload, &offset)
	if err != nil {
		return
	}

	if connectPkg.command != RTMP_AMF0_COMMAND_CONNECT {
		err = fmt.Errorf("handleAMF0CommandConnect command is error. command=", connectPkg.command)
		return
	}

	err, connectPkg.transactionId = Amf0ReadNumber(msg.payload, &offset)
	if err != nil {
		return
	}

	//this method is not strict for float type. just a warn.
	if 1 != connectPkg.transactionId {
		log.Println("warn:handleAMF0CommandConnect: transactionId is not 1. transactionId=", connectPkg.transactionId)
	}

	err, connectPkg.commandObjects = Amf0ReadObject(msg.payload, &offset)
	if err != nil {
		return
	}

	if offset < uint32(len(msg.payload)) {
		var v interface{}
		var marker uint8
		err, v, marker = Amf0Discovery(msg.payload, &offset)
		if err != nil {
			return
		}

		if RTMP_AMF0_Object == marker {
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

	return
}

func (pkg *Amf0CommandConnectPkg) Amf0ObjectsGetProperty(key string) (value interface{}) {

	for _, v := range pkg.commandObjects {
		if v.propertyName == key {
			return v.value
		}
	}

	return
}

func (rtmp *RtmpConn) ParseConnectPkg(pkg *Amf0CommandConnectPkg) (err error) {
	tcUrlValue := pkg.Amf0ObjectsGetProperty("tcUrl")
	if nil == tcUrlValue {
		err = fmt.Errorf("tcUrl is nil.")
		return
	}

	objectEncodingValue := pkg.Amf0ObjectsGetProperty("objectEncoding")
	if objectEncodingValue != nil {
		rtmp.objectEncoding = objectEncodingValue.(float64)
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
