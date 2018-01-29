package protocol

import (
	"encoding/binary"
	"fmt"
)

type SetChunkSizePacket struct {
	/**
	 * The maximum chunk size can be 65536 bytes. The chunk size is
	 * maintained independently for each direction.
	 */
	Chunk_size uint32
}

func (pkt *SetChunkSizePacket) Decode(data []uint8) (err error) {
	if len(data) < 4 {
		err = fmt.Errorf("decode set chunk size packet, data len is not enough.")
		return
	}

	var offset uint32
	pkt.Chunk_size = binary.BigEndian.Uint32(data[offset : offset+4])
	offset += 4

	return
}
func (pkt *SetChunkSizePacket) Encode() (data []uint8) {

	data = make([]uint8, 4)

	var offset uint32

	binary.BigEndian.PutUint32(data[offset:offset+4], pkt.Chunk_size)
	offset += 4

	return
}
func (pkt *SetChunkSizePacket) GetMessageType() uint8 {
	return RTMP_MSG_SetChunkSize
}
func (pkt *SetChunkSizePacket) GetPreferCsId() uint32 {
	return RTMP_CID_ProtocolControl
}
