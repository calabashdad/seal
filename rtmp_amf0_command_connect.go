package main

import (
	"fmt"
	"log"
)

type Amf0CommandConnectPkg struct {
	command        string
	transactionId  float64
	commandObjects []Amf0Object
	amfOptional    Amf0Object
}

func (rtmp *RtmpSession) handleAMF0CommandConnect(chunk *ChunkStruct) (err error) {

	var connectPkg Amf0CommandConnectPkg

	var offset uint32

	err, connectPkg.command = Amf0ReadString(chunk.msgPayload, &offset)
	if err != nil {
		return
	}

	if connectPkg.command != RTMP_AMF0_COMMAND_CONNECT {
		err = fmt.Errorf("handleAMF0CommandConnect command is error. command=", connectPkg.command)
		return
	}

	err, connectPkg.transactionId = Amf0ReadNumber(chunk.msgPayload, &offset)
	if err != nil {
		return
	}

	//this method is not strict for float type. just a warn.
	if 1 != connectPkg.transactionId {
		log.Println("warn:handleAMF0CommandConnect: transactionId is not 1. transactionId=", connectPkg.transactionId)
	}

	err, connectPkg.commandObjects = Amf0ReadObjects(chunk.msgPayload, &offset)
	if err != nil {
		return
	}

	if offset < uint32(len(chunk.msgPayload)) {
		var v interface{}
		var marker uint8
		err, v, marker = Amf0Discovery(chunk.msgPayload, &offset)
		if err != nil {
			return
		}

		if RTMP_AMF0_Object == marker {
			connectPkg.amfOptional = v.(Amf0Object)
		}
	}

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
