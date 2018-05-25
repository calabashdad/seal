package hls

import (
	"log"
	"seal/rtmp/pt"

	"github.com/calabashdad/utiltools"
)

// SourceStream delivery RTMP stream to HLS(m3u8 and ts),
type SourceStream struct {
	muxer *hlsMuxer
	cache *hlsCache

	codec  *avcAacCodec
	sample *codecSample
	jitter *pt.TimeJitter

	// we store the stream dts,
	// for when we notice the hls cache to publish,
	// it need to know the segment start dts.
	//
	// for example. when republish, the stream dts will
	// monotonically increase, and the ts dts should start
	// from current dts.
	//
	// or, simply because the HlsCache never free when unpublish,
	// so when publish or republish it must start at stream dts,
	// not zero dts.
	streamDts int64
}

// NewSourceStream new a hls source stream
func NewSourceStream() *SourceStream {
	return &SourceStream{
		muxer: newHlsMuxer(),
		cache: newHlsCache(),

		codec:  newAvcAacCodec(),
		sample: newCodecSample(),
		jitter: pt.NewTimeJitter(),
	}
}

// OnMeta process metadata
func (hls *SourceStream) OnMeta(pkt *pt.OnMetaDataPacket) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	if nil == pkt {
		return
	}

	if err = hls.codec.metaDataDemux(pkt); err != nil {
		return
	}

	return
}

// OnAudio process on audio data, mux to ts
func (hls *SourceStream) OnAudio(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	hls.sample.clear()
	if err = hls.codec.audioAacDemux(msg.Payload.Payload, hls.sample); err != nil {
		log.Println("hls codec demux audio failed, err=", err)
		return
	}

	if hls.codec.audioCodecID != pt.RtmpCodecAudioAAC {
		//log.Println("codec audio codec id is not aac, codeID=", hls.codec.audioCodecID)
		return
	}

	// ignore sequence header
	if pt.RtmpCodecAudioTypeSequenceHeader == hls.sample.aacPacketType {
		if err = hls.cache.onSequenceHeader(hls.muxer); err != nil {
			log.Println("hls cache on sequence header failed, err=", err)
			return
		}

		return
	}

	hls.jitter.Correct(msg, 0, 0, pt.RtmpTimeJitterFull)

	// the pts calc from rtmp/flv header
	pts := int64(msg.Header.Timestamp * 90)

	// for pure audio, update the stream dts also
	hls.streamDts = pts

	if err = hls.cache.writeAudio(hls.codec, hls.muxer, pts, hls.sample); err != nil {
		log.Println("hls cache write audio failed, err=", err)
		return
	}

	return
}

// OnVideo process on video data, mux to ts
func (hls *SourceStream) OnVideo(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	hls.sample.clear()
	if err = hls.codec.videoAvcDemux(msg.Payload.Payload, hls.sample); err != nil {
		log.Println("hls codec demuxer video failed, err=", err)
		return
	}

	// ignore info frame,
	if pt.RtmpCodecVideoAVCFrameVideoInfoFrame == hls.sample.frameType {
		return
	}

	if hls.codec.videoCodecID != pt.RtmpCodecVideoAVC {
		return
	}

	// ignore sequence header
	if pt.RtmpCodecVideoAVCFrameKeyFrame == hls.sample.frameType &&
		pt.RtmpCodecVideoAVCTypeSequenceHeader == hls.sample.frameType {
		return hls.cache.onSequenceHeader(hls.muxer)
	}

	hls.jitter.Correct(msg, 0, 0, pt.RtmpTimeJitterFull)

	dts := msg.Header.Timestamp * 90
	hls.streamDts = int64(dts)

	if err = hls.cache.writeVideo(hls.codec, hls.muxer, int64(dts), hls.sample); err != nil {
		log.Println("hls cache write video failed")
		return
	}

	return
}

// OnPublish publish stream event, continue to write the m3u8,
// for the muxer object not destroyed.
func (hls *SourceStream) OnPublish(app string, stream string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	if err = hls.cache.onPublish(hls.muxer, app, stream, hls.streamDts); err != nil {
		return
	}

	return
}

// OnUnPublish the unpublish event, only close the muxer, donot destroy the
// muxer, for when we continue to publish, the m3u8 will continue.
func (hls *SourceStream) OnUnPublish() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	if err = hls.cache.onUnPublish(hls.muxer); err != nil {
		return
	}

	return
}

func (hls *SourceStream) hlsMux() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()
	return
}
