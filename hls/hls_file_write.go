package hls

import (
	"os"
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

	return
}

func (fw *fileWriter) close() {

}

func (fw *fileWriter) isOpen() bool {
	return nil == fw.f
}

// write data to file
func (fw *fileWriter) write(buf []byte) (err error) {

	return
}
