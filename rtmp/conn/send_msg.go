package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/protocol"
)

func (rc *RtmpConn) SendMsg(msg *protocol.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	if nil == msg {
		return
	}

	// ensure the basic header is 1bytes. make simple
	if msg.Header.Perfer_csid < 2 {
		msg.Header.Perfer_csid = protocol.RTMP_CID_ProtocolControl
	}

	//current position of payload send.
	var offset uint32

	// always write the header event payload is empty.
	for {
		var header []uint8
	}

	return
}
