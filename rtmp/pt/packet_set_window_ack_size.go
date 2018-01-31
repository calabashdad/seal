package pt

import (
	"encoding/binary"
	"fmt"
)

/**
* 5.5. Window Acknowledgement Size (5)
* The client or the server sends this message to inform the peer which
* window size to use when sending acknowledgment.
 */
type SetWindowAckSizePacket struct {
	AckowledgementWindowSize uint32
}

func (pkt *SetWindowAckSizePacket) Decode(data []uint8) (err error) {
	if len(data) < 4 {
		err = fmt.Errorf("decode set window ack size packet, len is not enough.")
		return
	}

	var offset uint32
	pkt.AckowledgementWindowSize = binary.BigEndian.Uint32(data[offset : offset+4])
	offset += 4

	return
}
func (pkt *SetWindowAckSizePacket) Encode() (data []uint8) {

	data = make([]uint8, 4)

	var offset uint32
	binary.BigEndian.PutUint32(data[offset:offset+4], pkt.AckowledgementWindowSize)
	offset += 4

	return
}
func (pkt *SetWindowAckSizePacket) GetMessageType() uint8 {
	return RTMP_MSG_WindowAcknowledgementSize
}
func (pkt *SetWindowAckSizePacket) GetPreferCsId() uint32 {
	return RTMP_CID_ProtocolControl
}
