package pt

import (
	"encoding/binary"
)

// SetPeerBandWidthPacket The client or the server sends this message to update the output
// bandwidth of the peer.
type SetPeerBandWidthPacket struct {
	Bandwidth uint32
	TypeLimit uint8
}

// Decode .
func (pkt *SetPeerBandWidthPacket) Decode(data []uint8) (err error) {
	//nothing
	return
}

// Encode .
func (pkt *SetPeerBandWidthPacket) Encode() (data []uint8) {
	data = make([]uint8, 5)

	var offset uint32

	binary.BigEndian.PutUint32(data[offset:offset+4], pkt.Bandwidth)
	offset += 4

	data[offset] = pkt.TypeLimit
	offset++

	return
}

// GetMessageType .
func (pkt *SetPeerBandWidthPacket) GetMessageType() uint8 {
	return RTMP_MSG_SetPeerBandwidth
}

// GetPreferCsID .
func (pkt *SetPeerBandWidthPacket) GetPreferCsID() uint32 {
	return RtmpCidProtocolControl
}
