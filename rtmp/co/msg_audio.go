package co

import (
	"log"
	"seal/rtmp/flv"
	"seal/rtmp/pt"

	"github.com/calabashdad/utiltools"
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

	//copy to all consumers
	rc.source.copyToAllConsumers(msg)

	//cache the sequence.
	// do not cache the sequence header to gop cache, return here
	if flv.AudioIsSequenceHeader(msg.Payload.Payload) {
		rc.source.cacheAudioSequenceHeader = msg
		log.Println("cache audio data sequence")
		return
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

	return
}
