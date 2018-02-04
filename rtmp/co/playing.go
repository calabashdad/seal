package co

import (
	"fmt"
	"log"
	"seal/conf"
	"seal/rtmp/pt"
)

func (rc *RtmpConn) playing(p *pt.PlayPacket) (err error) {

	userSpecifiedDurationToStop := p.Duration > 0
	var startTime int64 = -1

	const timeOutUs = 1 * 1000 //ms

	for {
		//read from client. use short time out. 200 ms
		//if recv failed, it's ok, not an error.
		if true {
			var msg pt.Message
			_ = rc.RecvMsg(&msg.Header, &msg.Payload, timeOutUs)
			if len(msg.Payload.Payload) > 0 {
				//has recved play control.
				err = rc.handlePlayData(&msg)
				if err != nil {
					return
				}
			}
		}

		err, msgDump := rc.consumer.dump()
		if err != nil {
			break
		}

		// only when user specifies the duration,
		// we start to collect the durations for each message.
		if userSpecifiedDurationToStop {
			if startTime < 0 || startTime > int64(msgDump.Header.Timestamp) {
				startTime = int64(msgDump.Header.Timestamp)
			}

			rc.consumer.Duration += (float64(msgDump.Header.Timestamp) - float64(startTime))
			startTime = int64(msgDump.Header.Timestamp)
		}

		err = rc.SendMsg(msgDump, conf.GlobalConfInfo.Rtmp.TimeOut*1000000)
		if err != nil {
			break
		}

		log.Println("send msg to player, type=", msgDump.Header.MessageType, ",timestamp=", msgDump.Header.Timestamp,
			"msg payload=", len(msgDump.Payload.Payload))
	}

	return
}

func (rc *RtmpConn) handlePlayData(msg *pt.Message) (err error) {

	if nil == msg {
		return
	}

	if msg.Header.IsAmf0Data() || msg.Header.IsAmf3Data() {
		log.Println("play data: has not handled play data amf data. ")
		//todo.
	} else {
		//process user control
		rc.handlePlayUserControl(msg)
	}

	return
}

func (rc *RtmpConn) handlePlayUserControl(msg *pt.Message) (err error) {

	if nil == msg {
		return
	}

	if msg.Header.IsAmf0Command() || msg.Header.IsAmf3Command() {
		//ignore
		log.Println("ignore amf cmd when handle play user control.")
		return
	}

	// for jwplayer/flowplayer, which send close as pause message.
	if true {
		p := pt.CloseStreamPacket{}
		_ = p.Decode(msg.Payload.Payload)
		if 0 == len(p.CommandName) {
			//it's ok, not an error
			return
		} else {
			err = fmt.Errorf("player ask to close stream,remote=%s", rc.TcpConn.RemoteAddr())
			return
		}
	}

	// call msg,
	// support response null first
	if true {
		p := pt.CallPacket{}
		_ = p.Decode(msg.Payload.Payload)
		if 0 == len(p.CommandName) {
			//it's ok, not an error.
			log.Println("decode call packet failed when handle play user control.")
			return
		}

		pRes := pt.CallResPacket{}
		pRes.CommandName = pt.RTMP_AMF0_COMMAND_RESULT
		pRes.TransactionId = p.TransactionId
		pRes.CommandObjectMarker = pt.RTMP_AMF0_Null
		pRes.ResponseMarker = pt.RTMP_AMF0_Null

		err = rc.SendPacket(&pRes, 0, conf.GlobalConfInfo.Rtmp.TimeOut*1000000)
		if err != nil {
			return
		}
	}

	//pause or other msg
	if true {
		p := pt.PausePacket{}
		_ = p.Decode(msg.Payload.Payload)
		if 0 == len(p.CommandName) {
			//it's ok, not an error.
			log.Println("decode pause packet failed")
			return
		} else {
			err = rc.onPlayClientPause(uint32(rc.DefaultStreamId), p.IsPause)
			if err != nil {
				log.Println("play client pause error.err=", err)
				return
			}

			err = rc.consumer.onPlayPause(p.IsPause)
			if err != nil {
				log.Println("consumer on play pause error.err=", err)
				return
			}
		}
	}

	return
}

func (rc *RtmpConn) onPlayClientPause(streamId uint32, isPause bool) (err error) {

	if isPause {
		// onStatus(NetStream.Pause.Notify)
		p := pt.OnStatusCallPacket{}
		p.CommandName = pt.RTMP_AMF0_COMMAND_ON_STATUS

		p.Data = append(p.Data, pt.Amf0Object{
			PropertyName: pt.StatusLevel,
			Value:        pt.StatusLevelStatus,
			ValueType:    pt.RTMP_AMF0_String,
		})

		p.Data = append(p.Data, pt.Amf0Object{
			PropertyName: pt.StatusCode,
			Value:        pt.StatusCodeStreamPause,
			ValueType:    pt.RTMP_AMF0_String,
		})

		p.Data = append(p.Data, pt.Amf0Object{
			PropertyName: pt.StatusDescription,
			Value:        "Paused stream.",
			ValueType:    pt.RTMP_AMF0_String,
		})

		err = rc.SendPacket(&p, streamId, conf.GlobalConfInfo.Rtmp.TimeOut*1000000)
		if err != nil {
			return
		}

		// StreamEOF
		if true {
			p := pt.UserControlPacket{}
			p.EventType = pt.SrcPCUCStreamEOF
			p.EventData = streamId

			err = rc.SendPacket(&p, streamId, conf.GlobalConfInfo.Rtmp.TimeOut*1000000)
			if err != nil {
				log.Println("send PCUC(StreamEOF) message failed.")
				return
			}
		}

	} else {
		// onStatus(NetStream.Unpause.Notify)
		p := pt.OnStatusCallPacket{}
		p.CommandName = pt.RTMP_AMF0_COMMAND_ON_STATUS

		p.Data = append(p.Data, pt.Amf0Object{
			PropertyName: pt.StatusLevel,
			Value:        pt.StatusLevelStatus,
			ValueType:    pt.RTMP_AMF0_String,
		})

		p.Data = append(p.Data, pt.Amf0Object{
			PropertyName: pt.StatusCode,
			Value:        pt.StatusCodeStreamUnpause,
			ValueType:    pt.RTMP_AMF0_String,
		})

		p.Data = append(p.Data, pt.Amf0Object{
			PropertyName: pt.StatusDescription,
			Value:        "UnPaused stream.",
			ValueType:    pt.RTMP_AMF0_String,
		})

		err = rc.SendPacket(&p, streamId, conf.GlobalConfInfo.Rtmp.TimeOut*1000000)
		if err != nil {
			return
		}

		// StreanBegin
		if true {
			p := pt.UserControlPacket{}
			p.EventType = pt.SrcPCUCStreamBegin
			p.EventData = streamId

			err = rc.SendPacket(&p, streamId, conf.GlobalConfInfo.Rtmp.TimeOut*1000000)
			if err != nil {
				log.Println("send PCUC(StreanBegin) message failed.")
				return
			}
		}

	}

	return
}
