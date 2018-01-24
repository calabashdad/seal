package seal_rtmp_conn

import (
	"UtilsTools/identify_panic"
	"encoding/binary"
	"fmt"
	"log"
	"seal/seal_rtmp_server/seal_rtmp_protocol/protocol_stack"
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

	if chunkSize >= protocol_stack.RTMP_CHUNKSIZE_MIN && chunkSize <= protocol_stack.RTMP_CHUNKSIZE_MAX {
		rtmp.ChunkSize = chunkSize
		log.Println("peer set chunk size success. chunk size=", chunkSize)
	} else {
		//ignored
		log.Println("HandleMsgSetChunkSize, chunk size is invalid.", chunkSize)
	}

	return
}
