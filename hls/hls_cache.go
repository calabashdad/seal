package hls

import (
	"fmt"
	"log"
	"seal/conf"

	"github.com/calabashdad/utiltools"
	"seal/rtmp/pt"
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
	vf *mpegTsFrame
	vb []byte

	// the audio cache buffer start pts, to flush audio if full
	audioBufferStartPts int64
	// time jitter for aac
	aacJitter *hlsAacJitter
}

func newHlsCache() *hlsCache {
	return &hlsCache{
		af:        newMpegTsFrame(),
		vf:        newMpegTsFrame(),
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

	if err = hc.cacheVideo(codec, sample); err != nil {
		return
	}

	hc.vf.dts = dts
	hc.vf.pts = hc.vf.dts + int64(sample.cts)*int64(90)
	hc.vf.pid = tsVideoPid
	hc.vf.sid = tsVideoAvc
	hc.vf.key = sample.frameType == pt.RtmpCodecVideoAVCFrameKeyFrame

	// new segment when:
	// 1. base on gop.
	// 2. some gops duration overflow.
	if hc.vf.key && muxer.isSegmentOverflow() {
		if err = hc.reapSegment("video", muxer, hc.vf.dts); err != nil {
			return
		}
	}

	// flush video when got one
	if err = muxer.flushVideo(hc.af, hc.ab, hc.vf, &hc.vb); err != nil {
		log.Println("m3u8 muxer flush video failed")
		return
	}

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

	// for type1/5/6, insert aud packet.
	audNal := []byte{0x00, 0x00, 0x00, 0x01, 0x09, 0xf0}

	spsPpsSent := false
	audSent := false

	// a ts sample is format as:
	// 00 00 00 01 // header
	//       xxxxxxx // data bytes
	// 00 00 01 // continue header
	//       xxxxxxx // data bytes.
	// so, for each sample, we append header in aud_nal, then appends the bytes in sample.
	for i := 0; i < sample.nbSampleUnits; i++ {
		sampleUnit := sample.sampleUnits[i]
		size := len(sampleUnit.payload)

		if size <= 0 {
			return
		}

		// step 1:
		// first, before each "real" sample,
		// we add some packets according to the nal_unit_type,
		// for example, when got nal_unit_type=5, insert SPS/PPS before sample.

		// 5bits, 7.3.1 NAL unit syntax,
		// H.264-AVC-ISO_IEC_14496-10.pdf, page 44.
		var nalUnitType uint8
		nalUnitType = sampleUnit.payload[0]
		nalUnitType &= 0x1f

		// @see: ngx_rtmp_hls_video
		// Table 7-1 – NAL unit type codes, page 61
		// 1: Coded slice
		if 1 == nalUnitType {
			spsPpsSent = false
		}

		// 6: Supplemental enhancement information (SEI) sei_rbsp( ), page 61
		// @see: ngx_rtmp_hls_append_aud
		if !audSent {
			// @remark, when got type 9, we donot send aud_nal, but it will make ios unhappy, so we remove it.
			if 1 == nalUnitType || 5 == nalUnitType || 6 == nalUnitType {
				hc.vb = append(hc.vb, audNal...)
				audSent = true
			}
		}

		// 5: Coded slice of an IDR picture.
		// insert sps/pps before IDR or key frame is ok.
		if 5 == nalUnitType && !spsPpsSent {
			spsPpsSent = true

			// @see: ngx_rtmp_hls_append_sps_pps
			if codec.sequenceParameterSetLength > 0 {
				// AnnexB prefix, for sps always 4 bytes header
				hc.vb = append(hc.vb, audNal[:4]...)
				// sps
				hc.vb = append(hc.vb, codec.sequenceParameterSetNALUnit[:codec.sequenceParameterSetLength]...)
			}

			if codec.pictureParameterSetLength > 0 {
				// AnnexB prefix, for pps always 4 bytes header
				hc.vb = append(hc.vb, audNal[:4]...)
				// pps
				hc.vb = append(hc.vb, codec.pictureParameterSetNALUnit[:codec.pictureParameterSetLength]...)
			}
		}

		// 7-9, ignore, @see: ngx_rtmp_hls_video
		if nalUnitType >= 7 && nalUnitType <= 9 {
			continue
		}

		// step 2:
		// output the "real" sample, in buf.
		// when we output some special assist packets according to nal_unit_type

		// sample start prefix, '00 00 00 01' or '00 00 01'
		pAudnal := 0 + 1
		endAudnal := pAudnal + 3

		// first AnnexB prefix is long (4 bytes)
		if 0 == len(hc.vb) {
			pAudnal = 0
		}
		hc.vb = append(hc.vb, audNal[pAudnal:pAudnal+endAudnal-pAudnal]...)

		// sample data
		hc.vb = append(hc.vb, sampleUnit.payload...)
	}

	return
}
