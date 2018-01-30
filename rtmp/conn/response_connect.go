package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/pt"
)

func (rc *RtmpConn) ResponseConnect() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	var pkt pt.ConnectResPacket

	pkt.Command_name = pt.RTMP_AMF0_COMMAND_RESULT
	pkt.Transaction_id = 1

	pkt.Props = append(pkt.Props, pt.Amf0Object{
		PropertyName: "fmsVer",
		Value:        "FMS/" + pt.FMS_VERSION,
		ValueType:    pt.RTMP_AMF0_String,
	})

	pkt.Props = append(pkt.Props, pt.Amf0Object{
		PropertyName: "capabilities",
		Value:        127.0,
		ValueType:    pt.RTMP_AMF0_Number,
	})

	pkt.Props = append(pkt.Props, pt.Amf0Object{
		PropertyName: "mode",
		Value:        1.0,
		ValueType:    pt.RTMP_AMF0_Number,
	})

	pkt.Props = append(pkt.Props, pt.Amf0Object{
		PropertyName: pt.StatusLevel,
		Value:        pt.StatusLevelStatus,
		ValueType:    pt.RTMP_AMF0_String,
	})

	pkt.Props = append(pkt.Props, pt.Amf0Object{
		PropertyName: pt.StatusCode,
		Value:        pt.StatusCodeConnectSuccess,
		ValueType:    pt.RTMP_AMF0_String,
	})

	pkt.Props = append(pkt.Props, pt.Amf0Object{
		PropertyName: pt.StatusDescription,
		Value:        "Connection succeeded",
		ValueType:    pt.RTMP_AMF0_String,
	})

	pkt.Props = append(pkt.Props, pt.Amf0Object{
		PropertyName: "objectEncoding",
		Value:        rc.ConnectInfo.ObjectEncoding,
		ValueType:    pt.RTMP_AMF0_Number,
	})

	pkt.Info = append(pkt.Info, pt.Amf0Object{
		PropertyName: "version",
		Value:        pt.FMS_VERSION,
		ValueType:    pt.RTMP_AMF0_String,
	})

	pkt.Info = append(pkt.Info, pt.Amf0Object{
		PropertyName: "seal_license",
		Value:        "The MIT License (MIT)",
		ValueType:    pt.RTMP_AMF0_String,
	})

	pkt.Info = append(pkt.Info, pt.Amf0Object{
		PropertyName: "seal_authors",
		Value:        "YangKai",
		ValueType:    pt.RTMP_AMF0_String,
	})

	pkt.Info = append(pkt.Info, pt.Amf0Object{
		PropertyName: "seal_email",
		Value:        "beyondyangkai@gmail.com",
		ValueType:    pt.RTMP_AMF0_String,
	})

	pkt.Info = append(pkt.Info, pt.Amf0Object{
		PropertyName: "seal_copyright",
		Value:        "Copyright (c) 2018 YangKai",
		ValueType:    pt.RTMP_AMF0_String,
	})

	pkt.Info = append(pkt.Info, pt.Amf0Object{
		PropertyName: "seal_sig",
		Value:        "seal",
		ValueType:    pt.RTMP_AMF0_String,
	})

	err = rc.SendPacket(&pkt, 0)
	if err != nil {
		log.Println("response connect error.", err)
		return
	}

	return
}
