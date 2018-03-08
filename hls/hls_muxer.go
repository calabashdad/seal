package hls

import (
	"log"
	"os"
	"strconv"

	"github.com/calabashdad/utiltools"
)

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

// open a new segment, a new ts file
// segmentStartDts use to calc the segment duration, use 0 for the first segment of hls
func (hm *hlsMuxer) segmentOpen(segmentStartDts int64) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	if nil != hm.current {
		// has already opened, ignore segment open
		return
	}

	// create dir for app
	if err = hm.createDir(); err != nil {
		return
	}

	// new segment
	hm.current = newHlsSegment()
	hm.current.sequenceNo = hm.sequenceNo
	hm.sequenceNo++
	hm.current.segmentStartDts = segmentStartDts

	// generate filename
	filename := hm.stream + "-" + strconv.Itoa(hm.current.sequenceNo) + ".ts"

	hm.current.fullPath = hm.hlsPath + "/" + hm.app + "/" + filename
	hm.current.uri = filename

	tmpFile := hm.current.fullPath + ".tmp"
	if err = hm.current.muxer.open(tmpFile); err != nil {
		log.Println("open hls muxer failed, err=", err)
		return
	}

	return
}

func (hm *hlsMuxer) onSequenceHeader() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	return
}

// whether segment overflow,
// that is whether the current segment duration>=(the segment in config)
func (hm *hlsMuxer) isSegmentOverflow() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	return
}

// whether segment absolutely overflow, for pure audio to reap segment,
// that is whether the current segment duration>=2*(the segment in config)
func (hm *hlsMuxer) isSegmentAbsolutelyOverflow() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()
	return
}

func (hm *hlsMuxer) flushAudio(af *mpegTsFrame, ab []byte) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()
	return
}

func (hm *hlsMuxer) flushVideo(af *mpegTsFrame, ab []byte, vf *mpegTsFrame, vb []byte) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()
	return
}

// close segment(ts)
// logDesc is the description for log
func (hm *hlsMuxer) segmentClose(logDesc string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()
	return
}

func (hm *hlsMuxer) refreshM3u8() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()
	return
}

func (hm *hlsMuxer) _refreshM3u8() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()
	return
}

func (hm *hlsMuxer) createDir() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	appDir := hm.hlsPath
	appDir += "/"
	appDir += hm.app

	_, errLocal := os.Stat(appDir)
	if os.IsNotExist(errLocal) {
		if err = os.Mkdir(appDir, os.ModePerm); err != nil {
			return
		}
	}

	return
}
