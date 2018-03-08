package hls

import (
	"log"
	"seal/conf"

	"github.com/calabashdad/utiltools"
)

// hls stream cache,
// use to cache hls stream and flush to hls muxer.
//
// when write stream to ts file:
// video frame will directly flush to M3u8Muxer,
// audio frame need to cache, because it's small and flv tbn problem.
//
// whatever, the Hls cache used to cache video/audio,
// and flush video/audio to m3u8 muxer if needed.
//
// about the flv tbn problem:
//   flv tbn is 1/1000, ts tbn is 1/90000,
//   when timestamp convert to flv tbn, it will loose precise,
//   so we must gather audio frame together, and recalc the timestamp @see SrsHlsAacJitter,
//   we use a aac jitter to correct the audio pts.
type hlsCache struct {
	// current frame and buffer
	af *mpegTsFrame
	ab []byte
	vc *mpegTsFrame
	vb []byte

	// the audio cache buffer start pts, to flush audio if full
	audioBufferStartPts int64
	// time jitter for aac
	aacJitter *hlsAacJitter
}

func newHlsCache() *hlsCache {
	return &hlsCache{
		af:        newMpegTsFrame(),
		vc:        newMpegTsFrame(),
		aacJitter: newHlsAacJitter(),
	}
}

// when publish stream
func (hc *hlsCache) onPublish(muxer *hlsMuxer, app string, stream string, segmentStartDts int64) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	hlsFragment := conf.GlobalConfInfo.Hls.HlsFragment
	hlsWindow := conf.GlobalConfInfo.Hls.HlsWindow
	hlsPath := conf.GlobalConfInfo.Hls.HlsPath

	// open muxer
	if err = muxer.updateConfig(app, stream, hlsPath, hlsFragment, hlsWindow); err != nil {
		return
	}

	if err = muxer.segmentOpen(segmentStartDts); err != nil {
		return
	}

	return
}

// when unpublish stream
func (hc *hlsCache) onUnPublish(muxer *hlsMuxer) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	if err = muxer.flushAudio(hc.af, hc.ab); err != nil {
		log.Println("m3u8 muxer flush audio failed, err=", err)
		return
	}

	if err = muxer.segmentClose("unpublish"); err != nil {
		return
	}

	return
}

// when get sequence header,
// must write a #EXT-X-DISCONTINUITY to m3u8.
// @see: hls-m3u8-draft-pantos-http-live-streaming-12.txt
// @see: 3.4.11.  EXT-X-DISCONTINUITY
func (hc *hlsCache) onSequenceHeader(muxer *hlsMuxer) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()
	return
}

// write audio to cache, if need to flush, flush to muxer
func (hc *hlsCache) writeAudio(codec *avcAacCodec, muxer *hlsMuxer, pts int64, sample *codecSample) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()
	return
}

// wirte video to muxer
func (hc *hlsCache) writeVideo(codec *avcAacCodec, muxer *hlsMuxer, dts int64, sample *codecSample) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()
	return
}

// reopen the muxer for a new hls segment,
// close current segment, open a new segment,
// then write the key frame to the new segment.
// so, user must reap_segment then flush_video to hls muxer.
func (hc *hlsCache) reapSegment(logDesc string, muxer *hlsMuxer, segmentStartDts int64) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()
	return
}

func (hc *hlsCache) cacheAudio(codec *avcAacCodec, sample *codecSample) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()
	return
}

func (hc *hlsCache) cacheVideo(codec *avcAacCodec, sample *codecSample) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()
	return
}
