package main

import "log"

func (rtmp *RtmpSession) Connect() (err error) {

	var chunk *ChunkStruct
	err, chunk = rtmp.ExpectMsg()
	if err != nil {
		return
	}

	err = rtmp.DecodeMsg(chunk)
	if err != nil {
		return
	}

	connectPkg := chunk.decodeResult.(Amf0CommandConnectPkg)

	log.Println("rtmp connect result: ", connectPkg)

	return
}
