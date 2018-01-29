package conn

import (
	"seal/rtmp/protocol"
)

func (rc *RtmpConn) ResponsePingMsg(time_stamp uint32) (err error) {

	var pkt protocol.UserControlPacket

	pkt.Event_type = protocol.SrcPCUCPingResponse
	pkt.Event_data = time_stamp

	err = rc.SendPacket(&pkt, 0)
	if err != nil {
		return
	}

	return
}
