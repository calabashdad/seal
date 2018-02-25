package pt

// OnMetaDataPacket the stream metadata.
// @setDataFrame
type OnMetaDataPacket struct {

	// Name Name of metadata. Set to "onMetaData"
	Name string

	// Metadata Metadata of stream.
	Metadata interface{}

	// Marker type of Metadata object or ecma
	Marker uint8
}

// Decode .
func (pkt *OnMetaDataPacket) Decode(data []uint8) (err error) {
	var offset uint32

	if pkt.Name, err = Amf0ReadString(data, &offset); err != nil {
		return
	}

	if RTMP_AMF0_DATA_SET_DATAFRAME == pkt.Name {
		if pkt.Name, err = Amf0ReadString(data, &offset); err != nil {
			return
		}
	}

	if pkt.Metadata, err = amf0ReadAny(data, &pkt.Marker, &offset); err != nil {
		return
	}

	return
}

// Encode .
func (pkt *OnMetaDataPacket) Encode() (data []uint8) {
	data = append(data, amf0WriteString(pkt.Name)...)
	if RTMP_AMF0_Object == pkt.Marker {
		data = append(data, amf0WriteObject(pkt.Metadata.([]Amf0Object))...)
	} else if RTMP_AMF0_EcmaArray == pkt.Marker {
		data = append(data, amf0WriteEcmaArray(pkt.Metadata.(amf0EcmaArray))...)
	}

	return
}

// GetMessageType .
func (pkt *OnMetaDataPacket) GetMessageType() uint8 {
	return RtmpMsgAmf0DataMessage
}

// GetPreferCsID .
func (pkt *OnMetaDataPacket) GetPreferCsID() uint32 {
	return RtmpCidOverConnection2
}

// AddObject add object to objs
func (pkt *OnMetaDataPacket) AddObject(obj Amf0Object) {
	if RTMP_AMF0_Object == pkt.Marker {
		pkt.Metadata = append(pkt.Metadata.([]Amf0Object), obj)
	} else if RTMP_AMF0_EcmaArray == pkt.Marker {
		v := pkt.Metadata.(amf0EcmaArray)
		v.addObject(obj)

		pkt.Metadata = v
	}
}

// GetProperty get object property name
func (pkt *OnMetaDataPacket) GetProperty(name string) interface{} {

	if RTMP_AMF0_Object == pkt.Marker {
		for _, v := range pkt.Metadata.([]Amf0Object) {
			if name == v.propertyName {
				return v.value
			}
		}
	} else if RTMP_AMF0_EcmaArray == pkt.Marker {
		for _, v := range (pkt.Metadata.(amf0EcmaArray)).anyObject {
			if name == v.propertyName {
				return v.value
			}
		}
	}

	return nil
}
