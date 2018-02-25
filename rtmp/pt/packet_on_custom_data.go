package pt

// OnCustomDataPakcet  the stream custom data.
type OnCustomDataPakcet struct {

	// Name of custom data. Set to "onCustomData"
	Name string

	// Customdata Custom data of stream.
	Customdata interface{}

	// Marker type of CustomData
	Marker uint8
}

// Decode .
func (pkt *OnCustomDataPakcet) Decode(data []uint8) (err error) {
	var offset uint32

	if pkt.Name, err = Amf0ReadString(data, &offset); err != nil {
		return
	}

	if pkt.Customdata, err = amf0ReadAny(data, &pkt.Marker, &offset); err != nil {
		return
	}

	return
}

// Encode .
func (pkt *OnCustomDataPakcet) Encode() (data []uint8) {
	data = append(data, amf0WriteString(pkt.Name)...)
	if RTMP_AMF0_Object == pkt.Marker {
		data = append(data, amf0WriteObject(pkt.Customdata.([]Amf0Object))...)
	} else if RTMP_AMF0_EcmaArray == pkt.Marker {
		data = append(data, amf0WriteEcmaArray(pkt.Customdata.(amf0EcmaArray))...)
	}

	return
}

// GetMessageType .
func (pkt *OnCustomDataPakcet) GetMessageType() uint8 {
	return RtmpMsgAmf0DataMessage
}

// GetPreferCsID .
func (pkt *OnCustomDataPakcet) GetPreferCsID() uint32 {
	return RtmpCidOverConnection2
}
