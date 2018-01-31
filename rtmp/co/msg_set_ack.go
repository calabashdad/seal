package co

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/pt"
)

func (rc *RtmpConn) msgSetAck(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	log.Println("MsgSetChunk")

	p := pt.SetWindowAckSizePacket{}
	err = p.Decode(msg.Payload)
	if err != nil {
		return
	}

	if p.AckowledgementWindowSize > 0 {
		rc.AckWindow.AckWindowSize = p.AckowledgementWindowSize
		log.Println("set ack window size=", p.AckowledgementWindowSize)
	}

	return
}
