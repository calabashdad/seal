package seal_rtmp_conn

import (
	"time"
)

func (rtmp *RtmpConn) SendBytes(buf []uint8) (err error) {

	err = rtmp.Conn.SetDeadline(time.Now().Add(time.Duration(rtmp.TimeOut) * time.Second))
	if err != nil {
		return
	}

	var n int
	if n, err = rtmp.Conn.Write(buf); err != nil {
		return
	}
	_ = n

	return
}
