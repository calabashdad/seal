package hls

import (
	"log"
	"seal/rtmp/pt"

	"github.com/calabashdad/utiltools"
)

// the h264/avc and aac codec, for media stream.
//
// to demux the FLV/RTMP video/audio packet to sample,
// add each NALUs of h.264 as a sample unit to sample,
// while the entire aac raw data as a sample unit.
//
// for sequence header,
// demux it and save it in the avc_extra_data and aac_extra_data,
//
// for the codec info, such as audio sample rate,
// decode from FLV/RTMP header, then use codec info in sequence
// header to override it.
type avcAacCodec struct {
	// metadata specified
	duration  int
	width     int
	height    int
	frameRate int

	videoCodecID  int
	videoDataRate int // in bps
	audioCodecID  int
	audioDataRate int // in bps

	// video specified
	// profile_idc, H.264-AVC-ISO_IEC_14496-10.pdf, page 45.
	avcProfile uint8
	// level_idc, H.264-AVC-ISO_IEC_14496-10.pdf, page 45.
	avcLevel uint8
	// lengthSizeMinusOne, H.264-AVC-ISO_IEC_14496-15.pdf, page 16
	nalUnitLength               int8
	sequenceParameterSetLength  uint16
	sequenceParameterSetNALUnit []byte
	pictureParameterSetLength   uint16
	pictureParameterSetNALUnit  []byte

	// audio specified
	// 1.6.2.1 AudioSpecificConfig, in aac-mp4a-format-ISO_IEC_14496-3+2001.pdf, page 33.
	// audioObjectType, value defines in 7.1 Profiles, aac-iso-13818-7.pdf, page 40.
	aacProfile uint8
	// samplingFrequencyIndex
	aacSampleRate uint8
	// channelConfiguration
	aacChannels uint8

	// the avc extra data, the AVC sequence header,
	// without the flv codec header,
	// @see: ffmpeg, AVCodecContext::extradata
	avcExtraSize int
	avcExtraData []byte
	// the aac extra data, the AAC sequence header,
	// without the flv codec header,
	// @see: ffmpeg, AVCodecContext::extradata
	aacExtraSize int
	aacExtraData []byte
}

func newAvcAacCodec() *avcAacCodec {
	return &avcAacCodec{
		aacSampleRate: hlsAacSampleRateUnset,
	}
}

// demux the metadata, to get the stream info,
// for instance, the width/height, sample rate.
// @param metadata, the metadata amf0 object. assert not NULL.
func (codec *avcAacCodec) metaDataDemux(meta *pt.OnMetaDataPacket) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()
	return
}

// demux the audio packet in aac codec.
// the packet mux in FLV/RTMP format defined in flv specification.
// demux the audio speicified data(sound_format, sound_size, ...) to sample.
// demux the aac specified data(aac_profile, ...) to codec from sequence header.
// demux the aac raw to sample units.
func (codec *avcAacCodec) audioAacDemux(data []byte, sample *codecSample) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()
	return
}

// demux the video packet in h.264 codec.
// the packet mux in FLV/RTMP format defined in flv specification.
// demux the video specified data(frame_type, codec_id, ...) to sample.
// demux the h.264 sepcified data(avc_profile, ...) to codec from sequence header.
// demux the h.264 NALUs to sampe units.
func (codec *avcAacCodec) videoAvcDemux(data []byte, sample *codecSample) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()
	return
}
