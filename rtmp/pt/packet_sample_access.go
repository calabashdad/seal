package pt

// SampleAccessPacket .
type SampleAccessPacket struct {

	// CommandName Name of command. Set to "|RtmpSampleAccess".
	CommandName string

	// VideoSampleAccess whether allow access the sample of video.
	VideoSampleAccess bool

	// AudioSampleAccess whether allow access the sample of audio.
	AudioSampleAccess bool
}

// Decode .
func (pkt *SampleAccessPacket) Decode(data []uint8) (err error) {
	//nothing

	return
}

// Encode .
func (pkt *SampleAccessPacket) Encode() (data []uint8) {
	data = append(data, amf0WriteString(pkt.CommandName)...)
	data = append(data, amf0WriteBool(pkt.VideoSampleAccess)...)
	data = append(data, amf0WriteBool(pkt.AudioSampleAccess)...)

	return
}

// GetMessageType .
func (pkt *SampleAccessPacket) GetMessageType() uint8 {
	return RtmpMsgAmf0DataMessage
}

// GetPreferCsID .
func (pkt *SampleAccessPacket) GetPreferCsID() uint32 {
	return RtmpCidOverStream
}
