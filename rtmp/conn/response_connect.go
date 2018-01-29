package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/protocol"
)

func (rc *RtmpConn) ResponseConnect() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	var pkt protocol.ConnectResPacket

	pkt.Command_name = protocol.RTMP_AMF0_COMMAND_RESULT
	pkt.Transaction_id = 1

	pkt.Props = append(pkt.Props, protocol.Amf0Object{
		PropertyName: "fmsVer",
		Value:        "FMS/" + protocol.FMS_VERSION,
		ValueType:    protocol.RTMP_AMF0_String,
	})

	pkt.Props = append(pkt.Props, protocol.Amf0Object{
		PropertyName: "capabilities",
		Value:        127.0,
		ValueType:    protocol.RTMP_AMF0_Number,
	})

	pkt.Props = append(pkt.Props, protocol.Amf0Object{
		PropertyName: "mode",
		Value:        1.0,
		ValueType:    protocol.RTMP_AMF0_Number,
	})

	pkt.Props = append(pkt.Props, protocol.Amf0Object{
		PropertyName: protocol.StatusLevel,
		Value:        protocol.StatusLevelStatus,
		ValueType:    protocol.RTMP_AMF0_String,
	})

	pkt.Props = append(pkt.Props, protocol.Amf0Object{
		PropertyName: protocol.StatusCode,
		Value:        protocol.StatusCodeConnectSuccess,
		ValueType:    protocol.RTMP_AMF0_String,
	})

	pkt.Props = append(pkt.Props, protocol.Amf0Object{
		PropertyName: protocol.StatusDescription,
		Value:        "Connection succeeded",
		ValueType:    protocol.RTMP_AMF0_String,
	})

	pkt.Props = append(pkt.Props, protocol.Amf0Object{
		PropertyName: "objectEncoding",
		Value:        protocol.RTMP_SIG_AMF0_VER,
		ValueType:    protocol.RTMP_AMF0_Number,
	})

	pkt.Info = append(pkt.Info, protocol.Amf0Object{
		PropertyName: "version",
		Value:        protocol.FMS_VERSION,
		ValueType:    protocol.RTMP_AMF0_String,
	})

	pkt.Info = append(pkt.Info, protocol.Amf0Object{
		PropertyName: "seal_license",
		Value:        "The MIT License (MIT)",
		ValueType:    protocol.RTMP_AMF0_String,
	})

	pkt.Info = append(pkt.Info, protocol.Amf0Object{
		PropertyName: "seal_authors",
		Value:        "YangKai",
		ValueType:    protocol.RTMP_AMF0_String,
	})

	pkt.Info = append(pkt.Info, protocol.Amf0Object{
		PropertyName: "seal_email",
		Value:        "beyondyangkai@gmail.com",
		ValueType:    protocol.RTMP_AMF0_String,
	})

	pkt.Info = append(pkt.Info, protocol.Amf0Object{
		PropertyName: "seal_copyright",
		Value:        "Copyright (c) 2018 YangKai",
		ValueType:    protocol.RTMP_AMF0_String,
	})

	pkt.Info = append(pkt.Info, protocol.Amf0Object{
		PropertyName: "seal_sig",
		Value:        "seal",
		ValueType:    protocol.RTMP_AMF0_String,
	})

	err = rc.SendPacket(&pkt, 0)
	if err != nil {
		log.Println("response connect error.", err)
		return
	}

	return
}
