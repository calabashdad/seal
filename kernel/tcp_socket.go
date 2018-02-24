package kernel

import (
	"fmt"
	"io"
	"net"
	"time"
)

// TCPSock socket
type TCPSock struct {
	net.Conn
	recvTimeOut  uint32
	sendTimeOut  uint32
	recvBytesSum uint64
}

// GetRecvBytesSum return the recv bytes of conn
func (conn *TCPSock) GetRecvBytesSum() uint64 {
	return conn.recvBytesSum
}

// SetRecvTimeout set tcp socket recv timeout
func (conn *TCPSock) SetRecvTimeout(timeoutUs uint32) {
	conn.recvTimeOut = timeoutUs
}

// SetSendTimeout set tcp socket send timeout
func (conn *TCPSock) SetSendTimeout(timeoutUs uint32) {
	conn.sendTimeOut = timeoutUs
}

// ExpectBytesFull recv exactly size or error
func (conn *TCPSock) ExpectBytesFull(buf []uint8, size uint32) (err error) {

	if err = conn.SetDeadline(time.Now().Add(time.Duration(conn.recvTimeOut) * time.Microsecond)); err != nil {
		return
	}

	var n int
	if n, err = io.ReadFull(conn.Conn, buf[:size]); err != nil {
		return
	}

	conn.recvBytesSum += uint64(n)

	return
}

// SendBytes send buf
func (conn *TCPSock) SendBytes(buf []uint8) (err error) {

	if err = conn.SetDeadline(time.Now().Add(time.Duration(conn.sendTimeOut) * time.Microsecond)); err != nil {
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
