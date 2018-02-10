package co

import (
	"log"
	"seal/rtmp/flv"
	"seal/rtmp/pt"
	"utiltools"
)

func (rc *RtmpConn) msgVideo(msg *pt.Message) (err error) {
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

	//cache the key frame
	if flv.VideoH264IsSequenceHeaderAndKeyFrame(msg.Payload.Payload) {
		rc.source.cacheVideoSequenceHeader = msg
	}

	// rc.source.gopCache.cache(msg) todo.

	if rc.source.atc {
		if nil != rc.source.cacheAudioSequenceHeader {
			rc.source.cacheVideoSequenceHeader.Header.Timestamp = msg.Header.Timestamp
		}

		if nil != rc.source.cacheMetaData {
			rc.source.cacheMetaData.Header.Timestamp = msg.Header.Timestamp
		}
	}

	return
}
