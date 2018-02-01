package co

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/pt"
)

/**
* SendPacket the function has call the message
 */
func (rc *RtmpConn) SendPacket(pkt pt.Packet, streamID uint32) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
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

	err = rc.SendMsg(&msg)
	if err != nil {
		return
	}

	return
}