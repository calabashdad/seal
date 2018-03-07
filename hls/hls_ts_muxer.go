package hls

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
	return
}

func (tm *tsMuxer) writeAudio(af *mpegTsFrame, ab []byte) (err error) {
	return
}

func (tm *tsMuxer) writeVideo(vf *mpegTsFrame, vb []byte) (err error) {
	return
}

func (tm *tsMuxer) close() (err error) {

	return
}
