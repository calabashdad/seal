package co

import (
	"log"
	"seal/rtmp/pt"

	"github.com/calabashdad/utiltools"
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
	if err = p.Decode(msg.Payload.Payload); err != nil {
		return
	}

	if p.AckowledgementWindowSize > 0 {
		rc.ack.ackWindowSize = p.AckowledgementWindowSize
		log.Println("set ack window size=", p.AckowledgementWindowSize)
	}

	return
}
