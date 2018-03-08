package hls

import (
	"log"

	"github.com/calabashdad/utiltools"
)

// write data from frame(header info) and buffer(data) to ts file.
type tsMuxer struct {
	writer *fileWriter
	path   string
}

func newTsMuxer() *tsMuxer {
	return &tsMuxer{
		writer: newFileWriter(),
	}
}

func (tm *tsMuxer) open(path string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	tm.path = path

	tm.close()

	if err = tm.writer.open(tm.path); err != nil {
		log.Println("opem ts muxer path failed, err=", err)
		return
	}

	// write mpegts header
	if err = mpegtsWriteHeader(tm.writer); err != nil {
		log.Println("write mpegts header failed, err=", err)
		return
	}

	return
}

func (tm *tsMuxer) writeAudio(af *mpegTsFrame, ab []byte) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()
	return
}

func (tm *tsMuxer) writeVideo(vf *mpegTsFrame, vb []byte) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()
	return
}

func (tm *tsMuxer) close() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	tm.writer.close()

	return
}
