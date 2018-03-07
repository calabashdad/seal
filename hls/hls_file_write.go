package hls

import (
	"log"
	"os"

	"github.com/calabashdad/utiltools"
)

// write file
type fileWriter struct {
	file string
	f    *os.File
}

func newFileWriter() *fileWriter {
	return &fileWriter{}
}

func (fw *fileWriter) open(file string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()
	return
}

func (fw *fileWriter) close() {

}

func (fw *fileWriter) isOpen() bool {
	return nil == fw.f
}

// write data to file
func (fw *fileWriter) write(buf []byte) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()
	return
}
