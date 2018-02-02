package co

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/flv"
	"seal/rtmp/pt"
)

func (rc *RtmpConn) msgVideo(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	// log.Println("video data, csid=", msg.Header.PerferCsid,
	// 	",stream id=", msg.Header.StreamId,
	// 	", payload len=", len(msg.Payload.Payload),
	// 	",timestamp=", msg.Header.Timestamp)

	//copy to all consumers
	rc.source.copyToAllConsumers(msg)

	//cache the key frame
	if flv.VideoH264IsSequenceHeaderAndKeyFrame(msg.Payload.Payload) {
		rc.source.cacheVideoSequenceHeader = msg
	}

	rc.source.gopCache.cache(msg)

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
