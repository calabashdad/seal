package co

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/pt"
)

func (rc *RtmpConn) msgUserCtrl(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	log.Println("MsgUserCtrl")

	p := pt.UserControlPacket{}
	err = p.Decode(msg.Payload.Payload)
	if err != nil {
		return
	}

	switch p.EventType {
	case pt.SrcPCUCStreamBegin:
	case pt.SrcPCUCStreamEOF:
	case pt.SrcPCUCStreamDry:
	case pt.SrcPCUCSetBufferLength:
	case pt.SrcPCUCStreamIsRecorded:
	case pt.SrcPCUCPingRequest:
		err = rc.ctrlPingRequest(&p)
	case pt.SrcPCUCPingResponse:
	default:
		log.Println("msg user ctrl unknown event type.type=", p.EventType)
	}

	if err != nil {
		return
	}

	return
}

func (rc *RtmpConn) ctrlPingRequest(p *pt.UserControlPacket) (err error) {

	log.Println("ctrl ping request.")

	if pt.SrcPCUCSetBufferLength == p.EventType {

	} else if pt.SrcPCUCPingRequest == p.EventType {
		var pp pt.UserControlPacket
		pp.EventType = pt.SrcPCUCPingResponse
		pp.EventData = p.EventData
		err = rc.SendPacket(&pp, 0)
		if err != nil {
			return
		}

		log.Println("send ping response success.")

	}

	return
}
