package co

import (
	"log"
	"seal/rtmp/pt"
	"github.com/calabashdad/utiltools"
)

func (rc *RtmpConn) msgAggregate(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("aggregate")
	if nil == msg {
		return
	}

	return
}
