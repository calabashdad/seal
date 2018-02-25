package pt

import (
	"encoding/binary"
)

// AcknowlegementPacket The client or the server sends the acknowledgment to the peer after
// receiving bytes equal to the window size.
type AcknowlegementPacket struct {
	SequenceNumber uint32
}

// Decode .
func (pkt *AcknowlegementPacket) Decode(data []uint8) (err error) {
	//nothing
	return
}

// Encode .
func (pkt *AcknowlegementPacket) Encode() (data []uint8) {

	data = make([]uint8, 4)
	binary.BigEndian.PutUint32(data[:], pkt.SequenceNumber)

	return
}

// GetMessageType .
func (pkt *AcknowlegementPacket) GetMessageType() uint8 {
	return RTMP_MSG_Acknowledgement
}

// GetPreferCsID .
func (pkt *AcknowlegementPacket) GetPreferCsID() uint32 {
	return RtmpCidProtocolControl
}
