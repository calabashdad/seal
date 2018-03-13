package hls

import (
	"fmt"
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
		log.Println("segment open failed, err=", err)
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

	if err = muxer.flushAudio(hc.af, &hc.ab); err != nil {
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

	if err = muxer.onSequenceHeader(); err != nil {
		return
	}

	return
}

// write audio to cache, if need to flush, flush to muxer
func (hc *hlsCache) writeAudio(codec *avcAacCodec, muxer *hlsMuxer, pts int64, sample *codecSample) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	if 0 == len(hc.ab) {
		pts = hc.aacJitter.onBufferStart(pts, sample.soundRate, int(codec.aacSampleRate))

		hc.af.dts = pts
		hc.af.pts = pts
		hc.audioBufferStartPts = pts

		hc.af.pid = tsAudioPid
		hc.af.sid = tsAudioAac
	} else {
		hc.aacJitter.onBufferContinue()
	}

	// write audio to cache
	if err = hc.cacheAudio(codec, sample); err != nil {
		log.Println("hls cache audio failed, err=", err)
		return
	}

	if len(hc.ab) > hlsAudioCacheSize {
		if err = muxer.flushAudio(hc.af, &hc.ab); err != nil {
			log.Println("flush audio failed, err=", err)
			return
		}

	}

	//in ms, audio delay to flush the audios.
	var audioDelay = int64(hlsAacDelay)
	// flush if audio delay exceed
	if pts-hc.audioBufferStartPts > audioDelay*90 {
		if err = muxer.flushAudio(hc.af, &hc.ab); err != nil {
			return
		}
	}

	// reap when current source is pure audio.
	// it maybe changed when stream info changed,
	// for example, pure audio when start, audio/video when publishing,
	// pure audio again for audio disabled.
	// so we reap event when the audio incoming when segment overflow.
	// we use absolutely overflow of segment to make jwplayer/ffplay happy
	if muxer.isSegmentAbsolutelyOverflow() {
		if err = hc.reapSegment("audio", muxer, hc.af.pts); err != nil {
			log.Println("reap segment failed, err=", err)
			return
		}
		log.Println("reap segment success")
	}

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

	if err = muxer.segmentClose(logDesc); err != nil {
		log.Println("m3u8 muxer close segment failed, err=", err)
		return
	}

	if err = muxer.segmentOpen(segmentStartDts); err != nil {
		log.Println("m3u8 muxer open segment failed, err=", err)
		return
	}

	// segment open, flush the audio.
	// @see: ngx_rtmp_hls_open_fragment
	/* start fragment with audio to make iPhone happy */
	if err = muxer.flushAudio(hc.af, &hc.ab); err != nil {
		log.Println("m3u8 muxer flush audio failed, err=", err)
		return
	}

	return
}

func (hc *hlsCache) cacheAudio(codec *avcAacCodec, sample *codecSample) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	// AAC-ADTS
	// 6.2 Audio Data Transport Stream, ADTS
	// in aac-iso-13818-7.pdf, page 26.
	// fixed 7bytes header
	adtsHeader := [7]uint8{0xff, 0xf1, 0x00, 0x00, 0x00, 0x0f, 0xfc}

	for i := 0; i < sample.nbSampleUnits; i++ {
		sampleUnit := sample.sampleUnits[i]
		size := len(sampleUnit.payload)

		if size <= 0 || size > 0x1fff {
			err = fmt.Errorf("invalied aac frame length=%d", size)
			return
		}

		// the frame length is the AAC raw data plus the adts header size.
		frameLen := size + 7

		// adts_fixed_header
		// 2B, 16bits
		// int16_t syncword; //12bits, '1111 1111 1111'
		// int8_t ID; //1bit, '0'
		// int8_t layer; //2bits, '00'
		// int8_t protection_absent; //1bit, can be '1'

		// 12bits
		// int8_t profile; //2bit, 7.1 Profiles, page 40
		// TSAacSampleFrequency sampling_frequency_index; //4bits, Table 35, page 46
		// int8_t private_bit; //1bit, can be '0'
		// int8_t channel_configuration; //3bits, Table 8
		// int8_t original_or_copy; //1bit, can be '0'
		// int8_t home; //1bit, can be '0'

		// adts_variable_header
		// 28bits
		// int8_t copyright_identification_bit; //1bit, can be '0'
		// int8_t copyright_identification_start; //1bit, can be '0'
		// int16_t frame_length; //13bits
		// int16_t adts_buffer_fullness; //11bits, 7FF signals that the bitstream is a variable rate bitstream.
		// int8_t number_of_raw_data_blocks_in_frame; //2bits, 0 indicating 1 raw_data_block()

		// profile, 2bits
		adtsHeader[2] = (codec.aacProfile << 6) & 0xc0
		// sampling_frequency_index 4bits
		adtsHeader[2] |= (codec.aacSampleRate << 2) & 0x3c
		// channel_configuration 3bits
		adtsHeader[2] |= (codec.aacChannels >> 2) & 0x01
		adtsHeader[3] = (codec.aacChannels << 6) & 0xc0
		// frame_length 13bits
		adtsHeader[3] |= uint8((frameLen >> 11) & 0x03)
		adtsHeader[4] = uint8((frameLen >> 3) & 0xff)
		adtsHeader[5] = uint8((frameLen << 5) & 0xe0)
		// adts_buffer_fullness; //11bits
		adtsHeader[5] |= 0x1f

		// copy to audio buffer
		hc.ab = append(hc.ab, adtsHeader[:]...)
		hc.ab = append(hc.ab, sampleUnit.payload[:]...)

	}

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
