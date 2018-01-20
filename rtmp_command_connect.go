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

func (rtmp *RtmpSession) handleAMF0CommandConnect(chunk *ChunkStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}
	}()

	var connectPkg Amf0CommandConnectPkg

	var offset uint32

	err, connectPkg.command = Amf0ReadString(chunk.msg.payload, &offset)
	if err != nil {
		return
	}

	if connectPkg.command != RTMP_AMF0_COMMAND_CONNECT {
		err = fmt.Errorf("handleAMF0CommandConnect command is error. command=", connectPkg.command)
		return
	}

	err, connectPkg.transactionId = Amf0ReadNumber(chunk.msg.payload, &offset)
	if err != nil {
		return
	}

	//this method is not strict for float type. just a warn.
	if 1 != connectPkg.transactionId {
		log.Println("warn:handleAMF0CommandConnect: transactionId is not 1. transactionId=", connectPkg.transactionId)
	}

	err, connectPkg.commandObjects = Amf0ReadObject(chunk.msg.payload, &offset)
	if err != nil {
		return
	}

	if offset < uint32(len(chunk.msg.payload)) {
		var v interface{}
		var marker uint8
		err, v, marker = Amf0Discovery(chunk.msg.payload, &offset)
		if err != nil {
			return
		}

		if RTMP_AMF0_Object == marker {
			connectPkg.amfOptional = v
		}
	}

	chunk.decodeResultType = "Amf0CommandConnectPkg"
	chunk.decodeResult = connectPkg

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

func (rtmp *RtmpSession) ParseConnectPkg(pkg *Amf0CommandConnectPkg) (err error) {
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
