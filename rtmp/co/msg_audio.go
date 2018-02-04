package co

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/flv"
	"seal/rtmp/pt"
)

func (rc *RtmpConn) msgAudio(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	if nil == msg {
		return
	}

	// log.Println("audio data, csid=", msg.Header.PerferCsid,
	// 	",stream id=", msg.Header.StreamId,
	// 	", payload len=", len(msg.Payload.Payload),
	// 	",timestamp=", msg.Header.Timestamp)

	//cache the sequence.
	if flv.AudioIsSequenceHeader(msg.Payload.Payload) {
		rc.source.cacheAudioSequenceHeader = msg
	}

	rc.source.gopCache.cache(msg)

	if rc.source.atc {
		if nil != rc.source.cacheAudioSequenceHeader {
			rc.source.cacheAudioSequenceHeader.Header.Timestamp = msg.Header.Timestamp
		}

		if nil != rc.source.cacheMetaData {
			rc.source.cacheMetaData.Header.Timestamp = msg.Header.Timestamp
		}
	}

	//copy to all consumers
	rc.source.copyToAllConsumers(msg)

	return
}
