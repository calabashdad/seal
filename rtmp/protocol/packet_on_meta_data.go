package protocol

/**
* the stream metadata.
* FMLE: @setDataFrame
* others: onMetaData
 */

type OnMetaDataPacket struct {
}

func (pkt *OnMetaDataPacket) Decode([]uint8) (err error) {
	return
}
func (pkt *OnMetaDataPacket) Encode() (b []uint8) {
	return
}
func (pkt *OnMetaDataPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0DataMessage
}
func (pkt *OnMetaDataPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverConnection2
}
