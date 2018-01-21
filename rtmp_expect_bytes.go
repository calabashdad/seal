package main

import (
	"io"
	"time"
)

func (rtmp *RtmpConn) ExpectBytes(size uint32, buf []uint8) (err error) {

	err = rtmp.Conn.SetDeadline(time.Now().Add(time.Duration(g_conf_info.Rtmp.TimeOut) * time.Second))
	if err != nil {
		return
	}

	var recvSize int
	if recvSize, err = io.ReadFull(rtmp.Conn, buf[:size]); err != nil {
		return
	}

	rtmp.recvBytesSum += uint64(recvSize)

	return
}
