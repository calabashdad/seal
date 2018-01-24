package seal_rtmp_conn

import (
	"UtilsTools/identify_panic"
	"log"
)

func (rtmp *RtmpConn) handleUserCtrlSetBufferLength(chunkStreamId uint32, userCtrl *UserControlMsg) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	if err != nil {
		return
	}

	log.Println("handle user ctrl set buffer length ignore. user ctrl=", userCtrl)

	return
}
