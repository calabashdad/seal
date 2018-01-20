package main

import (
	"UtilsTools/identify_panic"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func (rtmp *RtmpSession) Connect() (err error) {

	var chunk *ChunkStream

	//expect connect msg.
	err, chunk = rtmp.RecvMsg()
	if err != nil {
		return
	}

	//decode connect msg
	err = rtmp.DecodeMsg(chunk)
	if err != nil {
		return
	}

	if "Amf0CommandConnectPkg" != chunk.decodeResultType {
		err = fmt.Errorf("can not expect connect message.")
		return
	}

	connectPkg := chunk.decodeResult.(Amf0CommandConnectPkg)
	log.Println("rtmp connect result: ", connectPkg)

	err = rtmp.ParseConnectPkg(&connectPkg)
	if err != nil {
		log.Println("parse connect pkg error.", err)
		return
	}

	err = rtmp.CommonMsgSetWindowAcknowledgementSize(chunk, 2500000)
	if err != nil {
		return
	}

	err = rtmp.CommonMsgSetPeerBandwidth(chunk, 2500000, 2)
	if err != nil {
		return
	}

	err = rtmp.ResponseConnectApp(chunk)
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
