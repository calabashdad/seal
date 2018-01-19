package main

import (
	"time"
)

func (rtmp *RtmpSession) SendBytes(buf []uint8) (err error) {

	err = rtmp.Conn.SetDeadline(time.Now().Add(time.Duration(g_conf_info.Rtmp.TimeOut) * time.Second))
	if err != nil {
		return
	}

	if _, err = rtmp.Conn.Write(buf); err != nil {
		return
	}

	return
}
