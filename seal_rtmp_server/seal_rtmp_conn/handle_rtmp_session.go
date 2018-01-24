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

		log.Println("One RtmpConn finished.remote=", rtmpSession.Conn.RemoteAddr())
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
	} else if RTMP_ROLE_PALY == rtmp.Role {
		rtmp.PlayerUnRegistePublishStream()
	}
}

func (rtmp *RtmpConn) PlayerRegistePublishStream() (res bool) {
	v, ok := MapPublishingStreams.Load(rtmp.StreamInfo.stream)
	if ok {
		rtmpPub := v.(*RtmpConn)
		rtmpPub.players.Store(rtmp.Conn.RemoteAddr(), rtmp.msgChan)
		res = true

		if nil != rtmpPub.cacheMsgMetaData {
			rtmp.msgChan <- rtmpPub.cacheMsgMetaData
		}

		if nil != rtmpPub.cacheMsgH264SequenceKeyFrame {
			rtmp.msgChan <- rtmpPub.cacheMsgH264SequenceKeyFrame
		}

		if nil != rtmpPub.cacheMsgAACSequenceHeader {
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

		log.Println("player dead, unregister from publisher.")
	}
}
