package co

import (
	"log"
	"seal/rtmp/pt"
	"github.com/calabashdad/utiltools"
)

func (rc *RtmpConn) msgAck(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("MsgAck")

	if nil == msg {
		return
	}

	p := pt.AcknowlegementPacket{}
	err = p.Decode(msg.Payload.Payload)

	return
}
