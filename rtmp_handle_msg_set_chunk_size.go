package main

import (
	"UtilsTools/identify_panic"
	"encoding/binary"
	"fmt"
	"log"
)

func (rtmp *RtmpConn) handleSetChunkSize(msg *MessageStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	var chunkSize uint32

	if len(msg.payload) >= 4 {
		chunkSize = binary.BigEndian.Uint32(msg.payload[0:4])
	} else {
		err = fmt.Errorf("handleSetChunkSize payload length < 4", len(msg.payload))
		return
	}

	if chunkSize >= RTMP_CHUNKSIZE_MIN && chunkSize <= RTMP_CHUNKSIZE_MAX {
		rtmp.chunkSize = chunkSize
		log.Println("peer set chunk size success. chunk size=", chunkSize)
	} else {
		//ignored
		log.Println("HandleMsgSetChunkSize, chunk size is invalid.", chunkSize)
	}

	return
}
