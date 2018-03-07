package hls

import (
	"log"

	"github.com/calabashdad/utiltools"
)

// jitter correct for audio,
// the sample rate 44100/32000 will lost precise,
// when mp4/ts(tbn=90000) covert to flv/rtmp(1000),
// so the Hls on ipad or iphone will corrupt,
// @see nginx-rtmp: est_pts
type hlsAacJitter struct {
	basePts   int64
	nbSamples int64
	syncMs    int
}

func newHlsAacJitter() *hlsAacJitter {
	return &hlsAacJitter{
		syncMs: hlsConfDefaultAacSync,
	}
}

// when buffer start, calc the "correct" pts for ts,
// @param flv_pts, the flv pts calc from flv header timestamp,
// @param sample_rate, the sample rate in format(flv/RTMP packet header).
// @param aac_sample_rate, the sample rate in codec(sequence header).
// @return the calc correct pts.
func (ha *hlsAacJitter) onBufferStart(flvPts int64, sampleRate int, aacSampleRate int) (calcCorrectPts int64) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	return
}

// when buffer continue, muxer donot write to file,
// the audio buffer continue grow and donot need a pts,
// for the ts audio PES packet only has one pts at the first time.
func (ha *hlsAacJitter) onBufferContinue() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	return
}
