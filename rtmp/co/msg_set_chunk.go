package co

import (
	"log"
	"seal/rtmp/pt"

	"github.com/calabashdad/utiltools"
)

func (rc *RtmpConn) msgSetChunk(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("set chunk size")

	if nil == msg {
		return
	}

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
