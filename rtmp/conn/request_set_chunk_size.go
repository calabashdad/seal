package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/pt"
)

func (rc *RtmpConn) RequestSetChunkSize(chunkSize uint32) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	var pkt pt.SetChunkSizePacket

	pkt.ChunkSize = chunkSize

	err = rc.SendPacket(&pkt, 0)
	if err != nil {
		return
	}

	return
}
