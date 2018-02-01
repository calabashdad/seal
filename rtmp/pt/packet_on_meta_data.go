package pt

/**
* the stream metadata.
* FMLE: @setDataFrame
* others: onMetaData
 */

type OnMetaDataPacket struct {
	/**
	 * Name of metadata. Set to "onMetaData"
	 */
	Name string
	/**
	 * Metadata of stream.
	 */
	Metadata interface{}
	Marker   uint8 //object or ecma
}

func (pkt *OnMetaDataPacket) Decode(data []uint8) (err error) {
	var offset uint32

	err, pkt.Name = Amf0ReadString(data, &offset)
	if err != nil {
		return
	}

	if RTMP_AMF0_DATA_SET_DATAFRAME == pkt.Name {
		err, pkt.Name = Amf0ReadString(data, &offset)
		if err != nil {
			return
		}
	}

	err, pkt.Metadata = Amf0ReadAny(data, &pkt.Marker, &offset)
	if err != nil {
		return
	}

	return
}
func (pkt *OnMetaDataPacket) Encode() (data []uint8) {
	data = append(data, Amf0WriteString(pkt.Name)...)
	if RTMP_AMF0_Object == pkt.Marker {
		data = append(data, Amf0WriteObject(pkt.Metadata.([]Amf0Object))...)
	} else if RTMP_AMF0_EcmaArray == pkt.Marker {
		data = append(data, Amf0WriteEcmaArray(pkt.Metadata.([]Amf0Object))...)
	}

	return
}
func (pkt *OnMetaDataPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0DataMessage
}
func (pkt *OnMetaDataPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverConnection2
}

func (pkt *OnMetaDataPacket) AddObject(obj Amf0Object) {
	if RTMP_AMF0_Object == pkt.Marker {
		pkt.Metadata = append(pkt.Metadata.([]Amf0Object), obj)
	} else if RTMP_AMF0_EcmaArray == pkt.Marker {
		v := pkt.Metadata.(Amf0EcmaArray)
		v.addObject(obj)

		pkt.Metadata = v
	}
}

func (pkt *OnMetaDataPacket) GetProperty(name string) interface{} {

	if RTMP_AMF0_Object == pkt.Marker {
		for _, v := range pkt.Metadata.([]Amf0Object) {
			if name == v.PropertyName {
				return v.Value
			}
		}
	} else if RTMP_AMF0_EcmaArray == pkt.Marker {
		for _, v := range (pkt.Metadata.(Amf0EcmaArray)).anyObject {
			if name == v.PropertyName {
				return v.Value
			}
		}
	}

	return nil
}
