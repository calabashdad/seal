package hls

// hlsMuxer the HLS stream(m3u8 and ts files).
// generally, the m3u8 muxer only provides methods to open/close segments,
// to flush video/audio, without any mechenisms.
//
// that is, user must use HlsCache, which will control the methods of muxer,
// and provides HLS mechenisms.
type hlsMuxer struct {
	app    string
	stream string

	hlsPath     string
	hlsFragment int
	hlsWindow   int

	sequenceNo int
	m3u8       string

	// m3u8 segments
	segments []*hlsSegment

	//current segment
	current *hlsSegment
}

func newHlsMuxer() *hlsMuxer {
	return &hlsMuxer{}
}

func (hm *hlsMuxer) getSequenceNo() int {
	return hm.sequenceNo
}

func (hm *hlsMuxer) updateConfig(app string, stream string, path string, fragment int, window int) (err error) {

	hm.app = app
	hm.stream = stream
	hm.hlsPath = path
	hm.hlsFragment = fragment
	hm.hlsWindow = window

	return
}

func (hm *hlsMuxer) segmentOpen(segmentStartDts int64) (err error) {

	if nil == hm.current {
		return
	}

	return
}

func (hm *hlsMuxer) onSequenceHeader() (err error) {

	return
}

// whether segment overflow,
// that is whether the current segment duration>=(the segment in config)
func (hm *hlsMuxer) isSegmentOverflow() (err error) {

	return
}

// whether segment absolutely overflow, for pure audio to reap segment,
// that is whether the current segment duration>=2*(the segment in config)
func (hm *hlsMuxer) isSegmentAbsolutelyOverflow() (err error) {

	return
}

func (hm *hlsMuxer) flushAudio(af *mpegTsFrame, ab []byte) (err error) {

	return
}

func (hm *hlsMuxer) flushVideo(af *mpegTsFrame, ab []byte, vf *mpegTsFrame, vb []byte) (err error) {

	return
}

// close segment(ts)
// logDesc is the description for log
func (hm *hlsMuxer) segmentClose(logDesc string) (err error) {
	return
}

func (hm *hlsMuxer) refreshM3u8() (err error) {
	return
}

func (hm *hlsMuxer) _refreshM3u8() (err error) {
	return
}

func (hm *hlsMuxer) createDir() (err error) {
	return
}
