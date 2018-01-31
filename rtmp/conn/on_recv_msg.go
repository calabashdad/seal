package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/pt"
)

func (rc *RtmpConn) OnRecvMsg(msg **pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	if rc.AckWindow.AckWindowSize > 0 &&
		((rc.TcpConn.RecvBytesSum - rc.AckWindow.HasAckedSize) > uint64(rc.AckWindow.AckWindowSize)) {
		//response a acknowlegement to peer.
		err = rc.ResponseAcknowlegementMsg()
		if err != nil {
			log.Println("response acknowlegement msg failed to peer.")
			return
		}
	}

	var pkt pt.Packet
	switch (*msg).Header.Message_type {
	case pt.RTMP_MSG_SetChunkSize, pt.RTMP_MSG_UserControlMessage, pt.RTMP_MSG_WindowAcknowledgementSize:

		err = rc.DecodeMsg(msg, &pkt)
		if err != nil {
			log.Println("decode msg faild. during on recv msg.")
			return
		}
	}

	switch (*msg).Header.Message_type {
	case pt.RTMP_MSG_SetChunkSize:
		chunkSize := pkt.(*pt.SetChunkSizePacket).ChunkSize
		if chunkSize >= pt.RTMP_CHUNKSIZE_MIN && chunkSize <= pt.RTMP_CHUNKSIZE_MAX {
			rc.InChunkSize = chunkSize

			log.Println("set in chunk size to ", chunkSize)
		}
	case pt.RTMP_MSG_UserControlMessage:
		if pt.SrcPCUCSetBufferLength == pkt.(*pt.UserControlPacket).Event_type {

		} else if pt.SrcPCUCPingRequest == pkt.(*pt.UserControlPacket).Event_type {
			err = rc.ResponsePingMsg(pkt.(*pt.UserControlPacket).Event_data)
		}
	case pt.RTMP_MSG_WindowAcknowledgementSize:
		ack_size := pkt.(*pt.SetWindowAckSizePacket).AckowledgementWindowSize
		if ack_size > 0 {
			rc.AckWindow.AckWindowSize = ack_size
		}
	}

	if err != nil {
		return
	}

	return
}
