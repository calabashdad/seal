package seal_rtmp_conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/seal_rtmp_server/seal_rtmp_protocol/protocol_stack"
)

func (rtmp *RtmpConn) handleMsgVideo(msg *MessageStream) (err error) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	//log.Println("recv video data, timestamp=", msg.header.timestamp)

	if rtmp.VideoIsH264SequenceAndKeyFrame(msg) {
		rtmp.cacheMsgH264SequenceKeyFrame = msg
	}

	//copy to all players
	rtmp.players.Range(func(key, value interface{}) bool {

		msgChanLocal := value.(chan *MessageStream)
		msgChanLocal <- msg

		log.Println("publisher put a video msg. timestamp=", msg.header.timestamp, ",msg payloadLen=", len(msg.payload))

		return true
	})

	return
}

func (rtmp *RtmpConn) VideoIsH264SequenceAndKeyFrame(msg *MessageStream) (res bool) {

	payloadLen := uint32(len(msg.payload))

	var offset uint32

	if payloadLen-offset < 2 {
		return false
	}

	codecId := msg.payload[0]
	codecId = codecId & 0x0f

	offset += 1

	if protocol_stack.SrsCodecVideoAVC != codecId {
		return false
	}

	frameType := msg.payload[0]
	frameType = (frameType >> 4) & 0x0F

	avcPacketType := msg.payload[1]

	if frameType == protocol_stack.SrsCodecVideoAVCFrameKeyFrame && avcPacketType == protocol_stack.SrsCodecVideoAVCTypeSequenceHeader {
		return true
	} else {
		return false
	}

	return
}
