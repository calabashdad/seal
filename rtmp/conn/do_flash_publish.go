package conn

import (
	"UtilsTools/identify_panic"
	"log"
)

func (rc *RtmpConn) DoFlashPublisherCycle() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	log.Println("new flash publisher come in. stream=", rc.StreamName)

	return
}
