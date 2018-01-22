package seal_rtmp_conn

import (
	"UtilsTools/identify_panic"
	"encoding/binary"
	"fmt"
	"log"
)

func (rtmp *RtmpConn) handleSetWindowAcknowledgementSize(msg *MessageStream) (err error) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	var windowAcknowlegementSize uint32
	if len(msg.payload) >= 4 {
		windowAcknowlegementSize = binary.BigEndian.Uint32(msg.payload[0:4])
	} else {
		err = fmt.Errorf("handleSetWindowAcknowledgementSize payload len < 4", len(msg.payload))
		return
	}

	if windowAcknowlegementSize > 0 {
		rtmp.AckWindow.ackWindowSize = windowAcknowlegementSize
		log.Println("peer set window acknowlegement size.", windowAcknowlegementSize)

	} else {
		//ignored.
		log.Println("HandleMsgSetWindowsAcknowlegementSize, ack size is invalied.", windowAcknowlegementSize)
	}

	return
}
