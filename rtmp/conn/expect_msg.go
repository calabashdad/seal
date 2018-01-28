package conn

import (
	"UtilsTools/identify_panic"
	"log"
)

//expect msg type, ignore other msgs until the type special has recv success.
func (rc *RtmpConn) ExpectMsg(packet interface{}) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	for {

		// rc.RecvMsg()

		// rc.DecodeMsg()
	}

	return
}
