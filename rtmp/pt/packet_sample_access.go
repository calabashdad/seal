package pt

type SampleAccessPacket struct {
	/**
	 * Name of command. Set to "|RtmpSampleAccess".
	 */
	CommandName string
	/**
	 * whether allow access the sample of video.
	 */
	VideoSampleAccess bool
	/**
	 * whether allow access the sample of audio.
	 */
	AudioSampleAccess bool
}

func (pkt *SampleAccessPacket) Decode(data []uint8) (err error) {
	//nothing

	return
}
func (pkt *SampleAccessPacket) Encode() (data []uint8) {
	data = append(data, Amf0WriteString(pkt.CommandName)...)
	data = append(data, Amf0WriteBool(pkt.VideoSampleAccess)...)
	data = append(data, Amf0WriteBool(pkt.AudioSampleAccess)...)

	return
}
func (pkt *SampleAccessPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0DataMessage
}
func (pkt *SampleAccessPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverStream
}
