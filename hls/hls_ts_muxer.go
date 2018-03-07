package hls

import (
	"log"

	"github.com/calabashdad/utiltools"
)

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
	return
}
