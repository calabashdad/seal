package pt

import (
	"encoding/binary"
)

type SetPeerBandWidthPacket struct {
	Bandwidth uint32
	TypeLimit uint8
}

func (pkt *SetPeerBandWidthPacket) Decode(data []uint8) (err error) {
	return
}
func (pkt *SetPeerBandWidthPacket) Encode() (data []uint8) {
	data = make([]uint8, 5)

	var offset uint32

	binary.BigEndian.PutUint32(data[offset:offset+4], pkt.Bandwidth)
	offset += 4

	data[offset] = pkt.TypeLimit
	offset++

	return
}
func (pkt *SetPeerBandWidthPacket) GetMessageType() uint8 {
	return RTMP_MSG_SetPeerBandwidth
}
func (pkt *SetPeerBandWidthPacket) GetPreferCsId() uint32 {
	return RTMP_CID_ProtocolControl
}
