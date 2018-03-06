package co

import (
	"log"
	"seal/rtmp/flv"
	"seal/rtmp/pt"

	"github.com/calabashdad/utiltools"
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

	// hls
	if nil != rc.source.hls {
		if err = rc.source.hls.onVideo(msg); err != nil {
			log.Println("hls process video data failed, err=", err)
			return
		}
	}

	//copy to all consumers
	rc.source.copyToAllConsumers(msg)

	//cache the key frame
	// do not cache the sequence header to gop cache, return here
	if flv.VideoH264IsSequenceHeaderAndKeyFrame(msg.Payload.Payload) {
		rc.source.cacheVideoSequenceHeader = msg
		log.Println("cache video sequence")
		return
	}

	rc.source.gopCache.cache(msg)

	if rc.source.atc {
		if nil != rc.source.cacheVideoSequenceHeader {
			rc.source.cacheVideoSequenceHeader.Header.Timestamp = msg.Header.Timestamp
		}

		if nil != rc.source.cacheMetaData {
			rc.source.cacheMetaData.Header.Timestamp = msg.Header.Timestamp
		}
	}

	return
}
