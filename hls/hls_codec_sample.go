package hls

import (
	"fmt"
	"log"
	"seal/rtmp/pt"

	"github.com/calabashdad/utiltools"
)

// the samples in the flv audio/video packet.
// the sample used to analysis a video/audio packet,
// split the h.264 NALUs to buffers, or aac raw data to a buffer,
// and decode the video/audio specified infos.
//
// the sample unit:
//       a video packet codec in h.264 contains many NALUs, each is a sample unit.
//       a audio packet codec in aac is a sample unit.
// @remark, the video/audio sequence header is not sample unit,
//       all sequence header stores as extra data,
// @remark, user must clear all samples before decode a new video/audio packet.
type codecSample struct {
	// each audio/video raw data packet will dumps to one or multiple buffers,
	// the buffers will write to hls and clear to reset.
	// generally, aac audio packet corresponding to one buffer,
	// where avc/h264 video packet may contains multiple buffer.
	nbSampleUnits int
	sampleUnits   [hlsMaxCodecSample]codecSampleUnit

	// whether the sample is video sample which demux from video packet.
	isVideo bool

	// CompositionTime, video_file_format_spec_v10_1.pdf, page 78.
	// cts = pts - dts, where dts = flvheader->timestamp.
	cts int

	// video specified
	frameType     int
	avcPacketType int

	// audio specified
	soundRate     int
	soundSize     int
	soundType     int
	aacPacketType int
}

func newCodecSample() *codecSample {
	return &codecSample{
		isVideo:       false,
		frameType:     pt.RtmpCodecVideoAVCFrameReserved,
		avcPacketType: pt.RtmpCodecVideoAVCTypeReserved,
		soundRate:     pt.RtmpCodecAudioSampleRateReserved,
		soundSize:     pt.RtmpCodecAudioSampleSizeReserved,
		soundType:     pt.RtmpCodecAudioSoundTypeReserved,
		aacPacketType: pt.RtmpCodecAudioTypeReserved,
	}
}

// clear all samples.
// in a word, user must clear sample before demux it.
func (sample *codecSample) clear() (err error) {

	sample.isVideo = false
	sample.nbSampleUnits = 0

	sample.cts = 0
	sample.frameType = pt.RtmpCodecVideoAVCFrameReserved
	sample.avcPacketType = pt.RtmpCodecVideoAVCTypeReserved

	sample.soundRate = pt.RtmpCodecAudioSampleRateReserved
	sample.soundSize = pt.RtmpCodecAudioSampleSizeReserved
	sample.soundType = pt.RtmpCodecAudioSoundTypeReserved
	sample.aacPacketType = pt.RtmpCodecAudioTypeReserved

	return
}

func (sample *codecSample) addSampleUnit(data []byte) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	if sample.nbSampleUnits >= hlsMaxCodecSample {
		err = fmt.Errorf("hls decode samples error, exceed the max count, nbSampleUnits=%d", sample.nbSampleUnits)
		return
	}

	sample.sampleUnits[sample.nbSampleUnits].payload = data
	sample.nbSampleUnits++

	return
}
