package hls

type mpegTsFrame struct {
	pts int64
	dts int64
	pid int
	sid int
	cc  int
	key bool
}

func newMpegTsFrame() *mpegTsFrame {
	return &mpegTsFrame{}
}
