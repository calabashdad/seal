package protocol

/**
* 4.1.3. createStream
* The client sends this command to the server to create a logical
* channel for message communication The publishing of audio, video, and
* metadata is carried out over stream channel created using the
* createStream command.
 */
type CreateStreamPacket struct {
}

func (pkt *CreateStreamPacket) Decode([]uint8) (err error) {
	return
}
func (pkt *CreateStreamPacket) Encode() (b []uint8) {
	return
}
func (pkt *CreateStreamPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}
func (pkt *CreateStreamPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverConnection
}
