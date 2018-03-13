package hls

// the wrapper of m3u8 segment from specification:
// 3.3.2.  EXTINF
// The EXTINF tag specifies the duration of a media segment.
type hlsSegment struct {
	// duration in seconds in m3u8.
	duration float64
	// sequence number in m3u8.
	sequenceNo int
	// ts uri in m3u8.
	uri string
	// ts full file to write.
	fullPath string
	// the muxer to write ts.
	muxer *tsMuxer
	// current segment start dts for m3u8
	segmentStartDts int64
	// whether current segement is sequence header.
	isSequenceHeader bool
}

func newHlsSegment() *hlsSegment {
	return &hlsSegment{
		muxer: newTsMuxer(),
	}
}

func (hs *hlsSegment) updateDuration(currentFrameDts int64) {

	// we use video/audio to update segment duration,
	// so when reap segment, some previous audio frame will
	// update the segment duration, which is nagetive,
	// just ignore it.
	if currentFrameDts < hs.segmentStartDts {
		return
	}

	hs.duration = float64(currentFrameDts-hs.segmentStartDts) / 90000.0
}
