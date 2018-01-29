package protocol

/**
* the special packet for the bandwidth test.
* actually, it's a OnStatusCallPacket, but
* 1. encode with data field, to send data to client.
* 2. decode ignore the data field, donot care.
 */
type BandWidthPacket struct {
}

func (pkt *BandWidthPacket) Decode([]uint8) (err error) {
	return
}

func (pkt *BandWidthPacket) Encode() (b []uint8) {
	return
}

func (pkt *BandWidthPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (pkt *BandWidthPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverStream
}
