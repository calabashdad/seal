package seal_rtmp_conn

import (
	"UtilsTools/identify_panic"
	"log"
)

func (rtmpSession *RtmpConn) HandleRtmpSession() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}

		rtmpSession.DestructWhenRtmpSessionDead()

		log.Println("One RtmpConn finished.remote=", rtmpSession.Conn.RemoteAddr(),
			",role=", rtmpSession.Role)
	}()

	log.Println("One RtmpConn come in. remote=", rtmpSession.Conn.RemoteAddr())

	err := rtmpSession.HandShake()
	if err != nil {
		log.Println("rtmp handshake failed, err=", err)
		return
	}

	log.Println("rtmp handshake success.remote=", rtmpSession.Conn.RemoteAddr())

	err = rtmpSession.RtmpMsgLoop()

	log.Println("rtmp msg loop quit.err=", err, ",remote=", rtmpSession.Conn.RemoteAddr())
}

func (rtmp *RtmpConn) DestructWhenRtmpSessionDead() {

	rtmp.Conn.Close()

	if RTMP_ROLE_PUBLISH == rtmp.Role {
		MapPublishingStreams.Delete(rtmp.StreamInfo.stream)
		log.Println("destruct, publisher deleted from MapPublishingStreams")
	} else if RTMP_ROLE_PALY == rtmp.Role {
		rtmp.PlayerUnRegistePublishStream()
		log.Println("player dead, unregister from publisher.")
	}
}

func (rtmp *RtmpConn) PlayerRegistePublishStream() (res bool) {
	v, ok := MapPublishingStreams.Load(rtmp.StreamInfo.stream)
	if ok {
		rtmpPub := v.(*RtmpConn)
		rtmpPub.players.Store(rtmp.Conn.RemoteAddr(), rtmp.msgChan)
		res = true

		if nil != rtmpPub.cacheMsgMetaData {
			log.Println("player, register, cache meta data.msg timestamp=", rtmpPub.cacheMsgMetaData.header.timestamp,
				",msg payloadSize=", len(rtmpPub.cacheMsgMetaData.payload))
			rtmp.msgChan <- rtmpPub.cacheMsgMetaData
		}

		if nil != rtmpPub.cacheMsgH264SequenceKeyFrame {
			log.Println("player, register, cache H264SequenceKeyFrame data.msg timestamp=", rtmpPub.cacheMsgH264SequenceKeyFrame.header.timestamp,
				",msg payloadSize=", len(rtmpPub.cacheMsgH264SequenceKeyFrame.payload))
			rtmp.msgChan <- rtmpPub.cacheMsgH264SequenceKeyFrame
		}

		if nil != rtmpPub.cacheMsgAACSequenceHeader {
			log.Println("player, register, cache AACSequenceHeader data.msg timestamp=", rtmpPub.cacheMsgAACSequenceHeader.header.timestamp,
				",msg payloadSize=", len(rtmpPub.cacheMsgAACSequenceHeader.payload))
			rtmp.msgChan <- rtmpPub.cacheMsgAACSequenceHeader
		}

	} else {
		res = false
	}

	return
}

func (rtmp *RtmpConn) PlayerUnRegistePublishStream() {
	v, ok := MapPublishingStreams.Load(rtmp.StreamInfo.stream)
	if ok {
		rtmpPub := v.(*RtmpConn)
		rtmpPub.players.Delete(rtmp.Conn.RemoteAddr())
	}
}
