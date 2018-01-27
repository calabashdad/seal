package seal_rtmp_conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/seal_rtmp_server/seal_rtmp_protocol/protocol_stack"
	"time"
)

func (rtmp *RtmpConn) handleMsgAudio(msg *MessageStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	//log.Println("recv audio data, timestamp=", msg.header.timestamp)

	if rtmp.AudioIsAACSequenceHeader(msg) {
		rtmp.cacheMsgAACSequenceHeader = msg
	}

	// copy to all players
	rtmp.players.Range(func(key, value interface{}) bool {

		msgChanLocal := value.(chan *MessageStream)

		select {
		case <-time.After(time.Millisecond * 100): //in case block
		case msgChanLocal <- msg:
			log.Println("publisher put a video msg, msg type=", msg.header.typeId)
		}

		return true
	})

	return
}

func (rtmp *RtmpConn) AudioIsAACSequenceHeader(msg *MessageStream) (res bool) {
	payloadLen := uint32(len(msg.payload))

	var offset uint32

	if payloadLen-offset < 2 {
		return false
	}

	soundFormat := msg.payload[0]
	soundFormat = (soundFormat >> 4) & 0x0f

	if soundFormat != protocol_stack.SrsCodecAudioAAC {
		return false
	}

	aac_packet_type := msg.payload[1]

	if aac_packet_type == protocol_stack.SrsCodecAudioTypeSequenceHeader {
		return true
	} else {
		return false
	}

	return
}
