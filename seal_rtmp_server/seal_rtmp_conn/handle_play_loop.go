package seal_rtmp_conn

import (
	"UtilsTools/identify_panic"
	"fmt"
	"log"
	"time"
)

func (rtmp *RtmpConn) handlePlayLoop() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	for {
		select {
		case <-time.After(time.Second * 5):
			err = fmt.Errorf("wait publisher put msg to chan timeout.")
			return
		case msg := <-rtmp.msgChan:

			log.Println("player pop a msg ")

			err = rtmp.SendMsg(msg)
			if err != nil {
				return
			}

		}

	}

	return
}
