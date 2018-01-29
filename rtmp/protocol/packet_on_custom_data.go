package protocol

/**
* the stream custom data.
* FMLE: @setDataFrame
* others: onCustomData
 */
type OnCustomDataPakcet struct {
}

func (pkt *OnCustomDataPakcet) Decode([]uint8) (err error) {
	return
}
func (pkt *OnCustomDataPakcet) Encode() (b []uint8) {
	return
}
func (pkt *OnCustomDataPakcet) GetMessageType() uint8 {
	return RTMP_MSG_AMF0DataMessage
}
func (pkt *OnCustomDataPakcet) GetPreferCsId() uint32 {
	return RTMP_CID_OverConnection2
}
