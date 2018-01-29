package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/protocol"
)

func (rc *RtmpConn) OnRecvMsg(msg *protocol.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	if rc.Ack_window.Ack_window_size > 0 &&
		((rc.TcpConn.RecvBytesSum - rc.Ack_window.Has_acked_size) > uint64(rc.Ack_window.Ack_window_size)) {
		//response a acknowlegement to peer.
		err = rc.ResponseAcknowlegementMsg()
		if err != nil {
			log.Println("response acknowlegement msg failed to peer.")
			return
		}
	}

	var pkt protocol.Packet
	switch msg.Header.Message_type {
	case protocol.RTMP_MSG_SetChunkSize, protocol.RTMP_MSG_UserControlMessage, protocol.RTMP_MSG_WindowAcknowledgementSize:

		err = rc.DecodeMsg(msg, pkt)
		if err != nil {
			log.Println("decode msg faild. during on recv msg.")
			return
		}
	}

	switch msg.Header.Message_type {
	case protocol.RTMP_MSG_SetChunkSize:
		chunk_size := pkt.(*protocol.SetChunkSizePacket).Chunk_size
		if chunk_size >= protocol.RTMP_CHUNKSIZE_MIN && chunk_size <= protocol.RTMP_CHUNKSIZE_MAX {
			rc.In_chunk_size = chunk_size
		}
	case protocol.RTMP_MSG_UserControlMessage:
		if protocol.SrcPCUCSetBufferLength == pkt.(*protocol.UserControlPacket).Event_type {

		} else if protocol.SrcPCUCPingRequest == pkt.(*protocol.UserControlPacket).Event_type {
			err = rc.ResponsePingMsg(pkt.(*protocol.UserControlPacket).Event_data)
		}
	case protocol.RTMP_MSG_WindowAcknowledgementSize:
		ack_size := pkt.(*protocol.SetWindowAckSizePacket).Ackowledgement_window_size
		if ack_size > 0 {
			rc.Ack_window.Ack_window_size = ack_size
		}
	}

	if err != nil {
		return
	}

	return
}
