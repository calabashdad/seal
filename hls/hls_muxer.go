package hls

import (
	"log"
	"os"
	"strconv"
	"syscall"

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
		log.Println("create dir faile,err=", err)
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

	// set the current segment to sequence header,
	// when close the segement, it will write a discontinuity to m3u8 file.
	hm.current.isSequenceHeader = true

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
func (hm *hlsMuxer) isSegmentAbsolutelyOverflow() bool {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	res := hm.current.duration >= float64(2*hm.hlsFragment)

	return res
}

func (hm *hlsMuxer) flushAudio(af *mpegTsFrame, ab *[]byte) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	// if current is NULL, segment is not open, ignore the flush event.
	if nil == hm.current {
		log.Println("hls segment is not open, ignore the flush event.")
		return
	}

	if len(*ab) <= 0 {
		return
	}

	hm.current.updateDuration(af.pts)

	if err = hm.current.muxer.writeAudio(af, *ab); err != nil {
		log.Println("current muxer write audio faile, err=", err)
		return
	}

	// write success, clear the buffer
	*ab = nil

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

	if nil == hm.current {
		return
	}

	// valid, add to segments if segment duration is ok
	if hm.current.duration*1000 >= hlsSegmentMinDurationMs {
		hm.segments = append(hm.segments, hm.current)

		// close the muxer of finished segment
		fullPath := hm.current.fullPath
		hm.current = nil

		// rename from tmp to real path
		tmpFile := fullPath + ".tmp"
		if err = os.Rename(tmpFile, fullPath); err != nil {
			return
		}
	} else {
		// reuse current segment index
		hm.sequenceNo--

		// rename from tmp to real path
		tmpFile := hm.current.fullPath + ".tmp"
		if err = syscall.Unlink(tmpFile); err != nil {
			log.Println("syscall unlink tmpfile=", tmpFile, " failed, err=", err)
		}
	}

	// the segment to remove
	var segmentToRemove []*hlsSegment

	// shrink the segments
	var duration float64
	removeIndex := -1
	for i := len(hm.segments) - 1; i >= 0; i-- {
		seg := hm.segments[i]
		duration += seg.duration

		if int(duration) > hm.hlsWindow {
			removeIndex = i
			break
		}
	}

	for i := 0; i < removeIndex; i++ {
		segmentToRemove = append(segmentToRemove, hm.segments[i])
	}
	hm.segments = hm.segments[removeIndex+1:]

	// refresh the m3u8, do not contains the removed ts
	hm.refreshM3u8()

	// remove the ts file
	for i := 0; i < len(segmentToRemove); i++ {
		s := segmentToRemove[i]
		syscall.Unlink(s.fullPath)
	}

	return
}

func (hm *hlsMuxer) refreshM3u8() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	m3u8File := hm.hlsPath
	m3u8File += "/"
	m3u8File += hm.app
	m3u8File += "/"
	m3u8File += hm.stream
	m3u8File += ".m3u8"

	hm.m3u8 = m3u8File
	m3u8File += ".temp"

	var f *os.File
	if f, err = hm._refreshM3u8(m3u8File); err != nil {
		log.Println("refresh m3u8 file faile, err=", err)
		return
	}
	if nil != f {
		f.Close()
		if err = os.Rename(m3u8File, hm.m3u8); err != nil {
			log.Println("rename m3u8 file failed, old file=", m3u8File, ",new file=", hm.m3u8)
			return
		}
	}

	syscall.Unlink(m3u8File)

	return
}

func (hm *hlsMuxer) _refreshM3u8(m3u8File string) (f *os.File, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	// no segments, return
	if 0 == len(hm.segments) {
		return
	}

	f, err = os.OpenFile(m3u8File, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("_refreshM3u8: open file error, file=", m3u8File)
		return
	}

	// #EXTM3U\n#EXT-X-VERSION:3\n
	header := []byte{
		// #EXTM3U\n
		0x23, 0x45, 0x58, 0x54, 0x4d, 0x33, 0x55, 0xa,
		// #EXT-X-VERSION:3\n
		0x23, 0x45, 0x58, 0x54, 0x2d, 0x58, 0x2d, 0x56, 0x45, 0x52,
		0x53, 0x49, 0x4f, 0x4e, 0x3a, 0x33, 0xa,
		// #EXT-X-ALLOW-CACHE:NO
		0x23, 0x45, 0x58, 0x54, 0x2d, 0x58, 0x2d, 0x41, 0x4c, 0x4c,
		0x4f, 0x57, 0x2d, 0x43, 0x41, 0x43, 0x48, 0x45, 0x3a, 0x4e, 0x4f, 0x0a,
	}

	if _, err = f.Write(header); err != nil {
		log.Println("write m3u8 header failed, err=", err)
		return
	}

	targetDuration := 0
	for i := 0; i < len(hm.segments); i++ {
		if int(hm.segments[i].duration) > targetDuration {
			targetDuration = int(hm.segments[i].duration)
		}
	}

	targetDuration++
	var duration string
	duration = "#EXT-X-TARGETDURATION:" + strconv.Itoa(targetDuration) + "\n"
	if _, err = f.Write([]byte(duration)); err != nil {
		log.Println("write m3u8 duration failed, err=", err)
		return
	}

	// write all segments
	for i := 0; i < len(hm.segments); i++ {
		s := hm.segments[i]

		if s.isSequenceHeader {
			// #EXT-X-DISCONTINUITY\n
			extDiscon := "#EXT-X-DISCONTINUITY\n"
			if _, err = f.Write([]byte(extDiscon)); err != nil {
				log.Println("write m3u8 segment discontinuity failed, err=", err)
				return
			}
		}

		// "#EXTINF:4294967295.208,\n"
		extInfo := "#EXTINF:" + strconv.FormatFloat(s.duration, 'f', 3, 64) + "\n"
		if _, err = f.Write([]byte(extInfo)); err != nil {
			log.Println("write m3u8 segment info failed, err=", err)
			return
		}

		// file name
		filename := s.uri + "\n"
		if _, err = f.Write([]byte(filename)); err != nil {
			log.Println("write m3u8 uri failed, err=", err)
			return
		}

	}

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
