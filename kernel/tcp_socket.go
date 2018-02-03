package kernel

import (
	"fmt"
	"io"
	"net"
	"time"
)

type TcpSock struct {
	net.Conn
	TimeOut      uint32
	RecvBytesSum uint64
}

func (conn *TcpSock) ExpectBytesFull(buf []uint8, size uint32, timeOutUs uint32) (err error) {

	err = conn.SetDeadline(time.Now().Add(time.Duration(timeOutUs) * time.Microsecond))
	if err != nil {
		return
	}

	var n int
	if n, err = io.ReadFull(conn.Conn, buf[:size]); err != nil {
		return
	}

	conn.RecvBytesSum += uint64(n)

	return
}

func (conn *TcpSock) SendBytes(buf []uint8, timeOutUs uint32) (err error) {

	err = conn.SetDeadline(time.Now().Add(time.Duration(timeOutUs) * time.Microsecond))
	if err != nil {
		return
	}

	var n int
	if n, err = conn.Conn.Write(buf); err != nil {
		return
	}

	if n != len(buf) {
		err = fmt.Errorf("tcp sock, send bytes error, need send ", len(buf), ",actually send ", n)
		return
	}

	return
}
