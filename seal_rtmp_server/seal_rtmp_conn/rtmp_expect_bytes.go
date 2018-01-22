package seal_rtmp_conn

import (
	"io"
	"time"
)

func (rtmp *RtmpConn) ExpectBytes(size uint32, buf []uint8) (err error) {

	err = rtmp.Conn.SetDeadline(time.Now().Add(time.Duration(rtmp.TimeOut) * time.Second))
	if err != nil {
		return
	}

	var recvSize int
	if recvSize, err = io.ReadFull(rtmp.Conn, buf[:size]); err != nil {
		return
	}

	rtmp.RecvBytesSum += uint64(recvSize)

	return
}
