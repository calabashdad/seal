package pt

import (
	"encoding/binary"
)

/**
* 5.3. Acknowledgement (3)
* The client or the server sends the acknowledgment to the peer after
* receiving bytes equal to the window size.
 */
type AcknowlegementPacket struct {
	SequenceNumber uint32
}

func (pkt *AcknowlegementPacket) Decode(data []uint8) (err error) {
	//nothing
	return
}

func (pkt *AcknowlegementPacket) Encode() (data []uint8) {

	data = make([]uint8, 4)
	binary.BigEndian.PutUint32(data[:], pkt.SequenceNumber)

	return
}

func (pkt *AcknowlegementPacket) GetMessageType() uint8 {
	return RTMP_MSG_Acknowledgement
}

func (pkt *AcknowlegementPacket) GetPreferCsId() uint32 {
	return RTMP_CID_ProtocolControl
}
