package conn

import (
	"seal/rtmp/protocol"
)

/**
* SendPacket the function has call the message 
*/
func (rc *RtmpConn) SendPacket(pkt protocol.Packet, streamID uint32) (err error) {

	var msg protocol.Message

	msg.Payload = pkt.Encode()
	msg.Size = uint32(len(msg.Payload))

	if uint32(len(msg.Payload)) <= 0 {
		//ignore empty msg.
		return
	}

	msg.Header.Payload_length = uint32(len(msg.Payload))
	msg.Header.Message_type = pkt.GetMessageType()
	msg.Header.Perfer_csid = pkt.GetPreferCsId()
	msg.Header.Stream_id = streamID

	err = rc.SendMsg(&msg)
	if err != nil {
		return
	}

	return
}
