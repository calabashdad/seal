package main

import (
	"io"
	"time"
)

func (rtmp *RtmpSession) ExpectBytes(size uint32, buf []uint8) (err error) {

	err = rtmp.Conn.SetDeadline(time.Now().Add(time.Duration(g_conf_info.Rtmp.TimeOut) * time.Second))
	if err != nil {
		return
	}

	if _, err = io.ReadFull(rtmp.Conn, buf[:size]); err != nil {
		return
	}

	return
}
