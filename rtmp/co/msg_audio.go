package co

import (
	"log"
	"seal/rtmp/flv"
	"seal/rtmp/pt"
	"utiltools"
)

func (rc *RtmpConn) msgAudio(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	if nil == msg {
		return
	}

	//cache the sequence.
	if flv.AudioIsSequenceHeader(msg.Payload.Payload) {
		rc.source.cacheAudioSequenceHeader = msg
	}

	// rc.source.gopCache.cache(msg) todo.

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
