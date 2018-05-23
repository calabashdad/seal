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

	// hls
	if nil != rc.source.hls {
		if err = rc.source.hls.OnAudio(msg); err != nil {
			log.Println("hls process audio data failed, err=", err)
			return
		}
	}

	//copy to all consumers
	rc.source.copyToAllConsumers(msg)

	//cache the sequence.
	// do not cache the sequence header to gop cache, return here
	if flv.AudioIsSequenceHeader(msg.Payload.Payload) {
		rc.source.CacheAudioSequenceHeader = msg
		log.Println("cache audio data sequence")
		return
	}

	rc.source.GopCache.cache(msg)

	//if rc.source.Atc {
	//	if nil != rc.source.CacheAudioSequenceHeader {
	//		rc.source.CacheAudioSequenceHeader.Header.Timestamp = msg.Header.Timestamp
	//	}
	//
	//	if nil != rc.source.CacheMetaData {
	//		rc.source.CacheMetaData.Header.Timestamp = msg.Header.Timestamp
	//	}
	//}

	return
}
