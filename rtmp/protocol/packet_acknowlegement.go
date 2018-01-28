package protocol

import (
	"encoding/binary"
)

/**
* 5.3. Acknowledgement (3)
* The client or the server sends the acknowledgment to the peer after
* receiving bytes equal to the window size.
 */
type AcknowlegementPacket struct {
	Sequence_number uint32
}

func (pkt *AcknowlegementPacket) Decode(b []uint8) (err error) {
	return
}

func (pkt *AcknowlegementPacket) Encode() (b []uint8) {

	b = make([]uint8, 4)
	binary.BigEndian.PutUint32(b[:], pkt.Sequence_number)

	return
}

func (pkt *AcknowlegementPacket) GetMessageType() uint8 {
	return RTMP_MSG_Acknowledgement
}

func (pkt *AcknowlegementPacket) GetPreferCsId() uint32 {
	return RTMP_CID_ProtocolControl
}
