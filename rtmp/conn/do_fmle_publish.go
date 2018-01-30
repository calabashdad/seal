package conn

import (
	"UtilsTools/identify_panic"
	"fmt"
	"log"
	"seal/rtmp/pt"
)

func (rc *RtmpConn) DoFmlePublisherCycle() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	log.Println("new fmle publisher come in. stream=", rc.StreamName)

	canPub := rc.CheckStreamCanPublish(rc.StreamName)
	if !canPub {
		err = fmt.Errorf("stream=%s can not publish, has already publishing.", rc.StreamName)
		return
	}

	//FCPublish
	var pktFmleStart *pt.FmleStartPacket
	var pktTmpFmleStart pt.Packet
	pktTmpFmleStart = pktFmleStart
	err = rc.ExpectMsg(&pktTmpFmleStart)
	if err != nil {
		log.Println("expect FCPublish msg failed.err=", err)
		return
	}
	pktFmleStart = pktTmpFmleStart.(*pt.FmleStartPacket)
	log.Println("expect fc publish packet success.")

	//FCPublish response
	var pktFmleStartRes pt.FmleStartResPacket
	pktFmleStartRes.Command_name = pt.RTMP_AMF0_COMMAND_RESULT
	pktFmleStartRes.Transaction_id = pktFmleStart.Transaction_id

	err = rc.SendPacket(&pktFmleStartRes, 0)
	if err != nil {
		log.Println("send fmle start res packet failed.err=", err)
		return
	}

	log.Println("fcpublish response success.")

	//createStream
	var pktCreateStream *pt.CreateStreamPacket
	var pktTmpCreateStream pt.Packet
	pktTmpCreateStream = pktCreateStream
	err = rc.ExpectMsg(&pktTmpCreateStream)
	if err != nil {
		log.Println("expect create stream failed.err=", err)
		return
	}
	pktCreateStream = pktTmpCreateStream.(*pt.CreateStreamPacket)
	log.Println("expect create stream packet success.")

	//createStream response
	var pktCreateStreamRes pt.CreateStreamResPacket
	pktCreateStreamRes.Command_name = pt.RTMP_AMF0_COMMAND_RESULT
	pktCreateStreamRes.Transaction_id = pktCreateStream.Transaction_id
	pktCreateStreamRes.Stream_id = rc.DefaultStreamId

	err = rc.SendPacket(&pktCreateStreamRes, 0)
	if err != nil {
		log.Println("send createStream response failed. err=", err)
		return
	}
	log.Println("send createStream response success.")

	//publish
	var pktPublish *pt.PublishPacket
	var pktTmpPublish pt.Packet
	pktTmpPublish = pktPublish
	err = rc.ExpectMsg(&pktTmpPublish)
	if err != nil {
		log.Println("expect publish packet faield.err=", err)
		return
	}
	pktPublish = pktTmpPublish.(*pt.PublishPacket)
	log.Println("expect publish packet success.")

	// publish response onFCPublish(NetStream.Publish.Start)
	var pktOnStatusCallonFCPublish pt.OnStatusCallPacket
	pktOnStatusCallonFCPublish.CommandName = pt.RTMP_AMF0_COMMAND_ON_FC_PUBLISH
	pktOnStatusCallonFCPublish.TransactionId = 0.0
	pktOnStatusCallonFCPublish.Data = append(pktOnStatusCallonFCPublish.Data, pt.Amf0Object{
		PropertyName: pt.StatusCode,
		Value:        pt.StatusCodePublishStart,
		ValueType:    pt.RTMP_AMF0_String,
	})

	pktOnStatusCallonFCPublish.Data = append(pktOnStatusCallonFCPublish.Data, pt.Amf0Object{
		PropertyName: pt.StatusDescription,
		Value:        "Started publishing stream, amazing",
		ValueType:    pt.RTMP_AMF0_String,
	})

	err = rc.SendPacket(&pktOnStatusCallonFCPublish, uint32(rc.DefaultStreamId))
	if err != nil {
		log.Println("send onFCPublish packet failed.err=", err)
		return
	}
	log.Println("send onFCPublish packet success.")

	// publish response onStatus(NetStream.Publish.Start)
	var pktOnStatusCallOnStatus pt.OnStatusCallPacket
	pktOnStatusCallOnStatus.CommandName = pt.RTMP_AMF0_COMMAND_ON_STATUS
	pktOnStatusCallOnStatus.TransactionId = 0.0 //for default
	pktOnStatusCallOnStatus.Data = append(pktOnStatusCallOnStatus.Data, pt.Amf0Object{
		PropertyName: pt.StatusLevel,
		Value:        pt.StatusLevelStatus,
		ValueType:    pt.RTMP_AMF0_String,
	})

	pktOnStatusCallOnStatus.Data = append(pktOnStatusCallOnStatus.Data, pt.Amf0Object{
		PropertyName: pt.StatusCode,
		Value:        pt.StatusCodePublishStart,
		ValueType:    pt.RTMP_AMF0_String,
	})

	pktOnStatusCallOnStatus.Data = append(pktOnStatusCallOnStatus.Data, pt.Amf0Object{
		PropertyName: pt.StatusDescription,
		Value:        "Started publishing stream. amazing.",
		ValueType:    pt.RTMP_AMF0_String,
	})

	pktOnStatusCallOnStatus.Data = append(pktOnStatusCallOnStatus.Data, pt.Amf0Object{
		PropertyName: pt.StatusClientId,
		Value:        pt.RTMP_SIG_CLIENT_ID,
		ValueType:    pt.RTMP_AMF0_String,
	})

	err = rc.SendPacket(&pktOnStatusCallOnStatus, uint32(rc.DefaultStreamId))
	if err != nil {
		log.Println("send on status packet failed.err=", err)
		return
	}
	log.Println("send on status packet success.")

	return
}
