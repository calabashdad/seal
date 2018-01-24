package seal_rtmp_conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/seal_rtmp_server/seal_rtmp_protocol/amf_serial"
	"seal/seal_rtmp_server/seal_rtmp_protocol/protocol_stack"
)

func (rtmp *RtmpConn) handleAmf0DataFrameOrMeta(msg *MessageStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	var offset uint32

	var metaDataName string
	err, metaDataName = amf_serial.Amf0ReadString(msg.payload, &offset)
	if err != nil {
		return
	}

	//ignore onMetaData.
	if protocol_stack.RTMP_AMF0_DATA_SET_DATAFRAME == metaDataName {
		err, metaDataName = amf_serial.Amf0ReadString(msg.payload, &offset)
		if err != nil {
			return
		}
	}

	err, rtmp.MetaData.value = amf_serial.Amf0ReadAny(msg.payload, &rtmp.MetaData.marker, &offset)

	if err != nil {
		return
	}

	if protocol_stack.RTMP_AMF0_EcmaArray == rtmp.MetaData.marker {
		log.Println("handle amf0 meta data success, meta data name=", metaDataName, ",ecma array=", rtmp.MetaData.value.(amf_serial.Amf0EcmaArray))

	} else if protocol_stack.RTMP_AMF0_Object == rtmp.MetaData.marker {
		log.Println("handle amf0 meta data success, meta data name=", metaDataName, ",obj=", rtmp.MetaData.value.([]amf_serial.Amf0Object))
	}

	return
}
