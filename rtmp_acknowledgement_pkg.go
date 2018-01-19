package main

import "encoding/binary"

func (rtmp *RtmpSession) AcknowledgementPkg() (err error, pkg []uint8) {

	//payload
	var payload [4]uint8

	sequenceNum := uint32(rtmp.recvBytesSum)
	binary.BigEndian.PutUint32(payload[:], sequenceNum)

	//header

	return
}
