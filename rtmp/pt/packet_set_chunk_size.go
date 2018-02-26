package pt

import (
	"encoding/binary"
	"fmt"
)

// SetChunkSizePacket Protocol control message 1, Set Chunk Size, is used to notify the
// peer about the new maximum chunk size.
type SetChunkSizePacket struct {
	// ChunkSize The maximum chunk size can be 65536 bytes. The chunk size is
	// maintained independently for each direction.
	ChunkSize uint32
}

// Decode .
func (pkt *SetChunkSizePacket) Decode(data []uint8) (err error) {
	if len(data) < 4 {
		err = fmt.Errorf("decode set chunk size packet, data len is not enough")
		return
	}

	var offset uint32
	pkt.ChunkSize = binary.BigEndian.Uint32(data[offset : offset+4])
	offset += 4

	return
}

// Encode .
func (pkt *SetChunkSizePacket) Encode() (data []uint8) {

	data = make([]uint8, 4)

	var offset uint32

	binary.BigEndian.PutUint32(data[offset:offset+4], pkt.ChunkSize)
	offset += 4

	return
}

// GetMessageType .
func (pkt *SetChunkSizePacket) GetMessageType() uint8 {
	return RtmpMsgSetChunkSize
}

// GetPreferCsID .
func (pkt *SetChunkSizePacket) GetPreferCsID() uint32 {
	return RtmpCidProtocolControl
}
