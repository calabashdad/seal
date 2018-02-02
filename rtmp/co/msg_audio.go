package co

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/pt"
)

func (rc *RtmpConn) msgAudio(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	log.Println("audio data, csid=", msg.Header.PerferCsid,
		",stream id=", msg.Header.StreamId,
		", payload len=", len(msg.Payload.Payload),
		",timestamp=", msg.Header.Timestamp)

	return
}
