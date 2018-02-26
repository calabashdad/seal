package co

import (
	"fmt"
	"log"
	"seal/conf"
	"seal/rtmp/pt"

	"github.com/calabashdad/utiltools"
)

func (rc *RtmpConn) msgAmf(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	if nil == msg {
		return
	}

	var offset uint32

	// skip 1bytes to decode the amf3 command.
	if pt.RtmpMsgAmf3CommandMessage == msg.Header.MessageType && msg.Header.PayloadLength >= 1 {
		offset++
	}

	//read the command name.
	var command string
	if command, err = pt.Amf0ReadString(msg.Payload.Payload, &offset); err != nil {
		return
	}

	if 0 == len(command) {
		err = fmt.Errorf("Amf0ReadString failed, command is nil.csid=%d", msg.Header.PerferCsid)
		return
	}

	log.Println("amf0/3 command or amf0/3 data, msg typeid=", msg.Header.MessageType, ",command =", command)

	switch command {
	case pt.RtmpAmf0CommandResult, pt.RtmpAmf0CommandError:
		err = rc.amf0ResultError(msg)
	case pt.RtmpAmf0CommandConnect:
		err = rc.amf0Connect(msg)
	case pt.RtmpAmf0CommandCreateStream:
		err = rc.amf0CreateStream(msg)
	case pt.RtmpAmf0CommandPlay:
		err = rc.amf0Play(msg)
	case pt.RtmpAmf0CommandPause:
		err = rc.amf0Pause(msg)
	case pt.RtmpAmf0CommandReleaseStream:
		err = rc.amf0ReleaseStream(msg)
	case pt.RtmpAmf0CommandFcPublish:
		err = rc.amf0FcPublish(msg)
	case pt.RtmpAmf0CommandPublish:
		err = rc.amf0Publish(msg)
	case pt.RtmpAmf0CommandUnpublish:
		err = rc.amf0UnPublish(msg)
	case pt.RtmpAmf0CommandKeeplive:
	case pt.RtmpAmf0CommandEnableVideo:
	case pt.RtmpAmf0DataSetDataFrame, pt.RtmpAmf0DataOnMetaData:
		err = rc.amf0Meta(msg)
	case pt.RtmpAmf0DataOnCustomData:
		err = rc.amf0OnCustomer(msg)
	case pt.RtmpAmf0CommandCloseStream:
		err = rc.amf0CloseStream(msg)
	case pt.RtmpAmf0CommandOnBwDone:
		err = rc.amf0OnBwDone(msg)
	case pt.RtmpAmf0CommandOnStatus:
		err = rc.amf0OnStatus(msg)
	case pt.RtmpAmf0CommandGetStreamLength:
		err = rc.amf0GetStreamLen(msg)
	case pt.RtmpAmf0CommandInsertKeyFrame:
	case pt.RtmpAmf0DataSampleAccess:
		err = rc.amf0SampleAccess(msg)
	default:
		log.Println("msg amf unknown command name=", command)
	}

	if err != nil {
		return
	}
	return
}

func (rc *RtmpConn) amf0ResultError(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("Amf0ResultError")

	if nil == msg {
		return
	}

	var offset uint32

	var transactionID float64
	if transactionID, err = pt.Amf0ReadNumber(msg.Payload.Payload, &offset); err != nil {
		log.Println("read transaction id failed when decode msg.", transactionID)
		return
	}

	reqCommandName := rc.cmdRequests[transactionID]
	if 0 == len(reqCommandName) {
		err = fmt.Errorf("can not find request command name-%s", reqCommandName)
		return
	}
	switch reqCommandName {
	case pt.RtmpAmf0CommandConnect:
		p := pt.ConnectResPacket{}
		err = p.Decode(msg.Payload.Payload)
	case pt.RtmpAmf0CommandCreateStream:
		p := pt.CreateStreamResPacket{}
		err = p.Decode(msg.Payload.Payload)
	case pt.RtmpAmf0CommandReleaseStream,
		pt.RtmpAmf0CommandFcPublish,
		pt.RtmpAmf0CommandUnpublish:
		p := pt.FmleStartResPacket{}
		err = p.Decode(msg.Payload.Payload)
	default:
		log.Println("result/error: unknown request command name=", reqCommandName)
	}

	if err != nil {
		log.Println("decode result or error response msg failed.")
		return
	}

	return
}

func (rc *RtmpConn) amf0Connect(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("Amf0Connect")

	if nil == msg {
		return
	}

	p := pt.ConnectPacket{}
	if err = p.Decode(msg.Payload.Payload); err != nil {
		log.Println("decode conncet pkt faile.err=", err)
		return
	}

	if nil == p.GetObjectProperty("tcUrl") {
		err = fmt.Errorf("no tcUrl info in connect")
		return
	}

	rc.connectInfo.tcURL = p.GetObjectProperty("tcUrl").(string)
	if o := p.GetObjectProperty("pageUrl"); o != nil {
		rc.connectInfo.pageURL = o.(string)
	}
	if o := p.GetObjectProperty("swfUrl"); o != nil {
		rc.connectInfo.swfURL = o.(string)
	}
	if o := p.GetObjectProperty("objectEncoding"); o != nil {
		rc.connectInfo.objectEncoding = o.(float64)
	}

	log.Println("decode connect pkt success.", p, ", info=", rc.connectInfo)

	var pkt pt.ConnectResPacket

	pkt.CommandName = pt.RtmpAmf0CommandResult
	pkt.TransactionID = 1

	pkt.AddProsObj(pt.NewAmf0Object("fmsVer", "FMS/"+pt.RtmpSigFmsVersion, pt.RtmpAmf0String))
	pkt.AddProsObj(pt.NewAmf0Object("capabilities", 127.0, pt.RtmpAmf0Number))
	pkt.AddProsObj(pt.NewAmf0Object("mode", 1.0, pt.RtmpAmf0Number))
	pkt.AddProsObj(pt.NewAmf0Object(pt.StatusLevel, pt.StatusLevelStatus, pt.RtmpAmf0String))
	pkt.AddProsObj(pt.NewAmf0Object(pt.StatusCode, pt.StatusCodeConnectSuccess, pt.RtmpAmf0String))
	pkt.AddProsObj(pt.NewAmf0Object(pt.StatusDescription, "Connection succeeded", pt.RtmpAmf0String))
	pkt.AddProsObj(pt.NewAmf0Object("objectEncoding", rc.connectInfo.objectEncoding, pt.RtmpAmf0Number))
	pkt.AddProsObj(pt.NewAmf0Object("version", pt.RtmpSigFmsVersion, pt.RtmpAmf0String))
	pkt.AddProsObj(pt.NewAmf0Object("seal_license", "The MIT License (MIT)", pt.RtmpAmf0String))
	pkt.AddProsObj(pt.NewAmf0Object("seal_authors", "YangKai", pt.RtmpAmf0String))
	pkt.AddProsObj(pt.NewAmf0Object("seal_email", "beyondyangkai@gmail.com", pt.RtmpAmf0String))
	pkt.AddProsObj(pt.NewAmf0Object("seal_copyright", "Copyright (c) 2018 YangKai", pt.RtmpAmf0String))
	pkt.AddProsObj(pt.NewAmf0Object("seal_sig", "seal", pt.RtmpAmf0String))

	if err = rc.SendPacket(&pkt, 0, conf.GlobalConfInfo.Rtmp.TimeOut*1000000); err != nil {
		log.Println("response connect error.", err)
		return
	}

	log.Println("send connect response success.")

	return
}

func (rc *RtmpConn) amf0CreateStream(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("Amf0CreateStream")

	if nil == msg {
		return
	}

	p := pt.CreateStreamPacket{}
	if err = p.Decode(msg.Payload.Payload); nil != err {
		log.Println("decode create stream failed.")
		return
	}

	//createStream response
	var pkt pt.CreateStreamResPacket
	pkt.CommandName = pt.RtmpAmf0CommandResult
	pkt.TransactionID = p.TransactionID
	pkt.StreamID = rc.defaultStreamID

	if err = rc.SendPacket(&pkt, 0, conf.GlobalConfInfo.Rtmp.TimeOut*1000000); err != nil {
		log.Println("send createStream response failed. err=", err)
		return
	}
	log.Println("send createStream response success.")

	return
}

func (rc *RtmpConn) amf0Play(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("Amf0Play")

	if nil == msg {
		return
	}

	p := pt.PlayPacket{}
	if err = p.Decode(msg.Payload.Payload); err != nil {
		return
	}

	log.Println("a new player come in, play info=", p)

	rc.streamName = p.StreamName

	source := gSources.findSourceToPlay(rc.streamName)
	if nil == source {
		err = fmt.Errorf("stream=%s can not play because has not published", rc.streamName)
		return
	}

	log.Println("play success. stream name=", rc.streamName)

	rc.source = source
	rc.role = RtmpRolePlayer

	//response start play.
	// StreamBegin
	if true {
		var pp pt.UserControlPacket
		pp.EventType = pt.SrcPCUCStreamBegin
		pp.EventData = msg.Header.StreamID
		err = rc.SendPacket(&pp, 0, conf.GlobalConfInfo.Rtmp.TimeOut*1000000)
		if err != nil {
			return
		}
		log.Println("send play stream begin pkt success.")
	}

	// onStatus(NetStream.Play.Reset)
	if true {
		var pp pt.OnStatusCallPacket
		pp.CommandName = pt.RtmpAmf0CommandOnStatus
		pp.AddObj(pt.NewAmf0Object(pt.StatusLevel, pt.StatusLevelStatus, pt.RtmpAmf0String))
		pp.AddObj(pt.NewAmf0Object(pt.StatusCode, pt.StatusCodeStreamReset, pt.RtmpAmf0String))
		pp.AddObj(pt.NewAmf0Object(pt.StatusDescription, "Playing and resetting stream.", pt.RtmpAmf0String))
		pp.AddObj(pt.NewAmf0Object(pt.StatusDetails, "stream", pt.RtmpAmf0String))
		pp.AddObj(pt.NewAmf0Object(pt.StatusClientID, pt.RtmpSigClientID, pt.RtmpAmf0String))

		if err = rc.SendPacket(&pp, uint32(rc.defaultStreamID), conf.GlobalConfInfo.Rtmp.TimeOut*1000000); err != nil {
			return
		}

		log.Println("send play onStatus(NetStream.Play.Reset) success.")
	}

	// onStatus(NetStream.Play.Start)
	if true {
		var pp pt.OnStatusCallPacket
		pp.CommandName = pt.RtmpAmf0CommandOnStatus

		pp.AddObj(pt.NewAmf0Object(pt.StatusLevel, pt.StatusLevelStatus, pt.RtmpAmf0String))
		pp.AddObj(pt.NewAmf0Object(pt.StatusCode, pt.StatusCodeStreamStart, pt.RtmpAmf0String))
		pp.AddObj(pt.NewAmf0Object(pt.StatusDescription, "Started playing stream.", pt.RtmpAmf0String))
		pp.AddObj(pt.NewAmf0Object(pt.StatusDetails, "stream", pt.RtmpAmf0String))
		pp.AddObj(pt.NewAmf0Object(pt.StatusClientID, pt.RtmpSigClientID, pt.RtmpAmf0String))

		if err = rc.SendPacket(&pp, uint32(rc.defaultStreamID), conf.GlobalConfInfo.Rtmp.TimeOut*1000000); err != nil {
			return
		}

		log.Println("send NetStream.Play.Reset response success.")
	}

	// |RtmpSampleAccess(false, false)
	if true {
		var pp pt.SampleAccessPacket
		pp.CommandName = pt.RtmpAmf0DataSampleAccess
		pp.AudioSampleAccess = true
		pp.VideoSampleAccess = true

		if err = rc.SendPacket(&pp, uint32(rc.defaultStreamID), conf.GlobalConfInfo.Rtmp.TimeOut*1000000); err != nil {
			return
		}

		log.Println("send RtmpSampleAccess success")
	}

	// onStatus(NetStream.Data.Start)
	if true {
		var pp pt.OnStatusDataPacket
		pp.CommandName = pt.RtmpAmf0CommandOnStatus

		pp.AddObj(pt.NewAmf0Object(pt.StatusLevel, pt.StatusLevelStatus, pt.RtmpAmf0String))
		pp.AddObj(pt.NewAmf0Object(pt.StatusCode, pt.StatusCodeDataStart, pt.RtmpAmf0String))
		pp.AddObj(pt.NewAmf0Object(pt.StatusDescription, "Started playing stream data.", pt.RtmpAmf0String))

		if err = rc.SendPacket(&pp, uint32(rc.defaultStreamID), conf.GlobalConfInfo.Rtmp.TimeOut*1000000); err != nil {
			return
		}

		log.Println("send NetStream.Data.Start success.")
	}

	rc.consumer = &Consumer{
		queueSizeMills: conf.GlobalConfInfo.Rtmp.ConsumerQueueSize * 1000,
		avStartTime:    -1,
		avEndTime:      -1,
		msgQuene:       make(chan *pt.Message, 1024),
	}

	rc.source.CreateConsumer(rc.consumer)

	if rc.source.atc && !rc.source.gopCache.empty() {
		if nil != rc.source.cacheMetaData {
			rc.source.cacheMetaData.Header.Timestamp = rc.source.gopCache.startTime()
		}
		if nil != rc.source.cacheVideoSequenceHeader {
			rc.source.cacheVideoSequenceHeader.Header.Timestamp = rc.source.gopCache.startTime()
		}
		if nil != rc.source.cacheAudioSequenceHeader {
			rc.source.cacheAudioSequenceHeader.Header.Timestamp = rc.source.gopCache.startTime()
		}
	}

	//cache meta data
	if nil != rc.source.cacheMetaData {
		rc.consumer.enquene(rc.source.cacheMetaData, rc.source.atc, rc.source.sampleRate, rc.source.frameRate, rc.source.timeJitter)
	}

	//cache video data
	if nil != rc.source.cacheVideoSequenceHeader {
		rc.consumer.enquene(rc.source.cacheVideoSequenceHeader, rc.source.atc, rc.source.sampleRate, rc.source.frameRate, rc.source.timeJitter)
	}

	//cache audio data
	if nil != rc.source.cacheAudioSequenceHeader {
		rc.consumer.enquene(rc.source.cacheAudioSequenceHeader, rc.source.atc, rc.source.sampleRate, rc.source.frameRate, rc.source.timeJitter)
	}

	//dump gop cache to client.
	rc.source.gopCache.dump(rc.consumer, rc.source.atc, rc.source.sampleRate, rc.source.frameRate, rc.source.timeJitter)

	log.Println("now playing, stream=", rc.streamName)

	err = rc.playing(&p)

	log.Println("playing over, stream=", rc.streamName)

	return
}

func (rc *RtmpConn) amf0Pause(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("Amf0Pause")

	if nil == msg {
		return
	}

	p := pt.PausePacket{}
	if err = p.Decode(msg.Payload.Payload); err != nil {
		return
	}

	return
}

func (rc *RtmpConn) amf0ReleaseStream(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("Amf0ReleaseStream")

	if nil == msg {
		return
	}

	p := pt.FmleStartPacket{}
	if err = p.Decode(msg.Payload.Payload); err != nil {
		return
	}

	var pp pt.FmleStartResPacket
	pp.CommandName = pt.RtmpAmf0CommandResult
	pp.TransactionID = p.TransactionID
	if err = rc.SendPacket(&pp, 0, conf.GlobalConfInfo.Rtmp.TimeOut*1000000); err != nil {
		return
	}
	log.Println("send release stream response success.")

	//set chunk size to peer.
	var pkt pt.SetChunkSizePacket
	pkt.ChunkSize = conf.GlobalConfInfo.Rtmp.ChunkSize
	if err = rc.SendPacket(&pkt, msg.Header.StreamID, conf.GlobalConfInfo.Rtmp.TimeOut*1000000); err != nil {
		return
	}
	log.Println("send request, set chunk size to ", pkt.ChunkSize)

	return
}

func (rc *RtmpConn) amf0FcPublish(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("Amf0FcPublish")

	if nil == msg {
		return
	}

	p := pt.FmleStartPacket{}
	if err = p.Decode(msg.Payload.Payload); err != nil {
		return
	}

	var pp pt.FmleStartResPacket
	pp.CommandName = pt.RtmpAmf0CommandResult
	pp.TransactionID = p.TransactionID
	if err = rc.SendPacket(&pp, 0, conf.GlobalConfInfo.Rtmp.TimeOut*1000000); err != nil {
		return
	}
	log.Println("send FcPublish response success.")

	return
}

func (rc *RtmpConn) amf0Publish(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("Amf0Publish")

	if nil == msg {
		return
	}

	p := pt.PublishPacket{}
	if err = p.Decode(msg.Payload.Payload); err != nil {
		return
	}

	log.Println("a new publisher come in, publish info=", p)

	rc.streamName = p.StreamName
	source := gSources.findSourceToPublish(rc.streamName)
	if nil == source {
		err = fmt.Errorf("stream=%s can not publish, find source is nil", rc.streamName)
		return
	}

	log.Println("stream=", rc.streamName, " published success.")

	rc.source = source
	rc.role = RtmpRoleFMLEPublisher

	var pp pt.OnStatusCallPacket
	pp.CommandName = pt.RtmpAmf0CommandOnStatus
	pp.AddObj(pt.NewAmf0Object(pt.StatusCode, pt.StatusCodePublishStart, pt.RtmpAmf0String))
	pp.AddObj(pt.NewAmf0Object(pt.StatusDescription, "Started publishing stream.", pt.RtmpAmf0String))

	if err = rc.SendPacket(&pp, uint32(rc.defaultStreamID), conf.GlobalConfInfo.Rtmp.TimeOut*1000000); err != nil {
		return
	}

	log.Println("send publish response success.")

	return
}

func (rc *RtmpConn) amf0UnPublish(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("Amf0UnPublish")

	if nil == msg {
		return
	}

	p := pt.FmleStartPacket{}
	if err = p.Decode(msg.Payload.Payload); err != nil {
		return
	}

	var pp pt.FmleStartResPacket
	pp.CommandName = pt.RtmpAmf0CommandResult
	pp.TransactionID = p.TransactionID
	if err = rc.SendPacket(&pp, 0, conf.GlobalConfInfo.Rtmp.TimeOut*1000000); err != nil {
		return
	}
	log.Println("send unpublish response success.")

	return
}

func (rc *RtmpConn) amf0Meta(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("Amf0Meta")

	if nil == msg {
		return
	}

	p := pt.OnMetaDataPacket{}
	if err = p.Decode(msg.Payload.Payload); err != nil {
		return
	}
	log.Println("decode meta data success, meta=", p)

	//add server info to metadata
	p.AddObject(*pt.NewAmf0Object("server", "seal rtmp server", pt.RtmpAmf0String))
	p.AddObject(*pt.NewAmf0Object("primary", "YangKai", pt.RtmpAmf0String))
	p.AddObject(*pt.NewAmf0Object("author", "YangKai", pt.RtmpAmf0String))

	if v := p.GetProperty("audiosamplerate"); v != nil {
		rc.source.sampleRate = v.(float64)
	}

	if v := p.GetProperty("framerate"); v != nil {
		rc.source.frameRate = v.(float64)
	}

	rc.source.atc = conf.GlobalConfInfo.Rtmp.Atc
	if v := p.GetProperty("bravo_atc"); v != nil {
		if conf.GlobalConfInfo.Rtmp.AtcAuto {
			rc.source.atc = true
		}
	}

	msg.Payload.Payload = p.Encode()
	msg.Header.PayloadLength = uint32(len(msg.Payload.Payload))

	// decode again and print
	if err = p.Decode(msg.Payload.Payload); err != nil {
		return
	}
	log.Println("meta data is ", p)

	//cache meta data
	if nil != rc.source {
		rc.source.cacheMetaData = msg
		log.Println("cache metadata")
	}

	rc.source.copyToAllConsumers(msg)

	return
}

func (rc *RtmpConn) amf0OnCustomer(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("Amf0OnCustomer")

	if nil == msg {
		return
	}

	p := pt.OnCustomDataPakcet{}
	if err = p.Decode(msg.Payload.Payload); err != nil {
		return
	}

	return
}

func (rc *RtmpConn) amf0CloseStream(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("Amf0CloseStream")

	if nil == msg {
		return
	}

	p := pt.CloseStreamPacket{}
	if err = p.Decode(msg.Payload.Payload); err != nil {
		return
	}

	return
}

func (rc *RtmpConn) amf0OnBwDone(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("Amf0OnBwDone")

	if nil == msg {
		return
	}

	return
}

func (rc *RtmpConn) amf0OnStatus(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("Amf0Onstats")

	if nil == msg {
		return
	}

	return
}

func (rc *RtmpConn) amf0GetStreamLen(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("Amf0GetStreamLen")

	if nil == msg {
		return
	}

	return
}

func (rc *RtmpConn) amf0SampleAccess(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("Amf0SampleAccess")

	if nil == msg {
		return
	}

	return
}
