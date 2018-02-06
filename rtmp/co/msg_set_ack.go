package co

import (
	"log"
	"seal/rtmp/pt"
	"utiltools"
)

func (rc *RtmpConn) msgSetAck(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("MsgSetChunk")

	if nil == msg {
		return
	}

	p := pt.SetWindowAckSizePacket{}
	err = p.Decode(msg.Payload.Payload)
	if err != nil {
		return
	}

	if p.AckowledgementWindowSize > 0 {
		rc.AckWindow.AckWindowSize = p.AckowledgementWindowSize
		log.Println("set ack window size=", p.AckowledgementWindowSize)
	}

	return
}
