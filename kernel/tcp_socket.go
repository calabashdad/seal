package kernel

import (
	"fmt"
	"io"
	"net"
	"time"
)

// TcpSock socket
type TcpSock struct {
	net.Conn
	recvTimeOut  uint32
	sendTimeOut  uint32
	RecvBytesSum uint64
}

func (conn *TcpSock) SetRecvTimeout(timeoutUs uint32) {
	conn.recvTimeOut = timeoutUs
}

func (conn *TcpSock) SetSendTimeout(timeoutUs uint32) {
	conn.sendTimeOut = timeoutUs
}

func (conn *TcpSock) ExpectBytesFull(buf []uint8, size uint32) (err error) {

	err = conn.SetDeadline(time.Now().Add(time.Duration(conn.recvTimeOut) * time.Microsecond))
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

func (conn *TcpSock) SendBytes(buf []uint8) (err error) {

	err = conn.SetDeadline(time.Now().Add(time.Duration(conn.sendTimeOut) * time.Microsecond))
	if err != nil {
		return
	}

	var n int
	if n, err = conn.Conn.Write(buf); err != nil {
		return
	}

	if n != len(buf) {
		err = fmt.Errorf("tcp sock, send bytes error, need send %d, actually send %d", len(buf), n)
		return
	}

	return
}
