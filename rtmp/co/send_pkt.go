package co

import (
	"log"
	"seal/rtmp/pt"

	"github.com/calabashdad/utiltools"
)

/**
* SendPacket the function has call the message
 */
func (rc *RtmpConn) SendPacket(pkt pt.Packet, streamID uint32, timeOutUs uint32) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	var msg pt.Message

	msg.Payload.Payload = pkt.Encode()
	msg.Payload.SizeTmp = uint32(len(msg.Payload.Payload))

	if uint32(len(msg.Payload.Payload)) <= 0 {
		//ignore empty msg.
		return
	}

	msg.Header.PayloadLength = uint32(len(msg.Payload.Payload))
	msg.Header.MessageType = pkt.GetMessageType()
	msg.Header.PerferCsid = pkt.GetPreferCsId()
	msg.Header.StreamId = streamID

	err = rc.SendMsg(&msg, timeOutUs)
	if err != nil {
		return
	}

	return
}
