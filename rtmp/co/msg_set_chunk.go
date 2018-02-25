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
	if err = p.Decode(msg.Payload.Payload); err != nil {
		return
	}

	if p.ChunkSize >= pt.RtmpChunkSizeMin && p.ChunkSize <= pt.RtmpChunkSizeMax {
		rc.inChunkSize = p.ChunkSize
		log.Println("peer set chunk size to ", p.ChunkSize)
	}

	return
}
