package co

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/pt"
)

func (rc *RtmpConn) msgSetChunk(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	log.Println("set chunk size")

	p := pt.SetChunkSizePacket{}
	err = p.Decode(msg.Payload.Payload)
	if err != nil {
		return
	}

	if p.ChunkSize >= pt.RTMP_CHUNKSIZE_MIN && p.ChunkSize <= pt.RTMP_CHUNKSIZE_MAX {
		rc.InChunkSize = p.ChunkSize
		log.Println("peer set chunk size to ", p.ChunkSize)
	}

	return
}
