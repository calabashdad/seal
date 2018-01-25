package seal_rtmp_conn

import (
	"UtilsTools/identify_panic"
	"encoding/binary"
	"log"
	"seal/seal_rtmp_server/seal_rtmp_protocol/protocol_stack"
)

func (rtmp *RtmpConn) SetChunkSize(chunkStreamId uint32, chunkSize uint32) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	var msg MessageStream

	//msg payload
	msg.payload = make([]uint8, 4)
	binary.BigEndian.PutUint32(msg.payload, chunkSize)

	//msg header
	msg.header.length = uint32(len(msg.payload))
	msg.header.typeId = protocol_stack.RTMP_MSG_SetChunkSize
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
