package pt

/**
* the stream custom data.
* FMLE: @setDataFrame
* others: onCustomData
 */
type OnCustomDataPakcet struct {
	/**
	* Name of custom data. Set to "onCustomData"
	 */
	Name string
	/**
	* Custom data of stream.
	 */
	Customdata interface{}
	Marker     uint8
}

func (pkt *OnCustomDataPakcet) Decode(data []uint8) (err error) {
	var offset uint32

	pkt.Name, err = Amf0ReadString(data, &offset)
	if err != nil {
		return
	}

	pkt.Customdata, err = Amf0ReadAny(data, &pkt.Marker, &offset)
	if err != nil {
		return
	}

	return
}
func (pkt *OnCustomDataPakcet) Encode() (data []uint8) {
	data = append(data, Amf0WriteString(pkt.Name)...)
	if RTMP_AMF0_Object == pkt.Marker {
		data = append(data, Amf0WriteObject(pkt.Customdata.([]Amf0Object))...)
	} else if RTMP_AMF0_EcmaArray == pkt.Marker {
		data = append(data, Amf0WriteEcmaArray(pkt.Customdata.(Amf0EcmaArray))...)
	}

	return
}
func (pkt *OnCustomDataPakcet) GetMessageType() uint8 {
	return RTMP_MSG_AMF0DataMessage
}
func (pkt *OnCustomDataPakcet) GetPreferCsId() uint32 {
	return RTMP_CID_OverConnection2
}
