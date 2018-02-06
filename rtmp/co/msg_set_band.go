package co

import (
	"log"
	"seal/rtmp/pt"
	"utiltools"
)

func (rc *RtmpConn) msgSetBand(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("MsgSetBand")

	if nil == msg {
		return
	}

	p := pt.SetPeerBandWidthPacket{}
	err = p.Decode(msg.Payload.Payload)

	return
}
