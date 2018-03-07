package hls

import (
	"seal/rtmp/pt"
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

	return
}

// OnAudio process on audio data, mux to ts
func (hls *SourceStream) OnAudio(msg *pt.Message) (err error) {

	return
}

// OnVideo process on video data, mux to ts
func (hls *SourceStream) OnVideo(msg *pt.Message) (err error) {

	return
}

// OnPublish publish stream event, continue to write the m3u8,
// for the muxer object not destroyed.
func (hls *SourceStream) OnPublish(app string, stream string) (err error) {
	if err = hls.cache.onPublish(hls.muxer, app, stream, hls.streamDts); err != nil {
		return
	}

	return
}

// the unpublish event, only close the muxer, donot destroy the
// muxer, for when we continue to publish, the m3u8 will continue.
func (hls *SourceStream) onUnPublish() (err error) {

	return
}

func (hls *SourceStream) hlsMux() {

	return
}
