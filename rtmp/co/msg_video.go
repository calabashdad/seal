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

	//if flv.VideoH264IsKeyframe(msg.Payload.Payload) {
	//	log.Printf("+++++recv i frame msg, type=%d, time=%d, len=%d\n", msg.Header.MessageType, msg.Header.Timestamp, msg.Header.PayloadLength)
	//}

	// hls
	if nil != rc.source.hls {
		if err = rc.source.hls.OnVideo(msg); err != nil {
			log.Println("hls process video data failed, err=", err)
			return
		}
	}

	//copy to all consumers
	rc.source.copyToAllConsumers(msg)

	//cache the key frame
	// do not cache the sequence header to gop cache, return here
	if flv.VideoH264IsKeyFrameAndSequenceHeader(msg.Payload.Payload) {
		rc.source.CacheVideoSequenceHeader = msg
		log.Println("cache video sequence")
		return
	}

	rc.source.GopCache.cache(msg)

	if rc.source.Atc {
		if nil != rc.source.CacheVideoSequenceHeader {
			rc.source.CacheVideoSequenceHeader.Header.Timestamp = msg.Header.Timestamp
		}

		if nil != rc.source.CacheMetaData {
			rc.source.CacheMetaData.Header.Timestamp = msg.Header.Timestamp
		}
	}

	return
}
