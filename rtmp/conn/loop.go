package conn

import (
	"UtilsTools/identify_panic"
	"log"
)

func (rc *RtmpConn) Loop() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	var err error

	err = rc.SetWindowAckSize(2500000)
	if err != nil {
		log.Println("set window ack size error.", err)
		return
	}

	err = rc.SetPeerBandWidth(2500000, 2)
	if err != nil {
		log.Println("set peer band width error.", err)
		return
	}

	//todo. bandwidth test.

	err = rc.ResponseConnect()
	if err != nil {
		return
	}

	err = rc.OnBwDone()
	if err != nil {
		return
	}

	err = rc.IdentifyClient()
	if err != nil {
		log.Println("identify client type failed.err=", err)
		return
	}

}
