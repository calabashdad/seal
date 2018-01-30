package pt

//message header.
type MessageHeader struct {
	/**
	 * 3bytes.
	 * Three-byte field that contains a timestamp delta of the message.
	 * The 4 bytes are packed in the big-endian order.
	 * @remark, only used for decoding message from chunk stream.
	 */
	Timestamp_delta uint32
	/**
	 * 3bytes.
	 * Three-byte field that represents the size of the payload in bytes.
	 * It is set in big-endian format.
	 */
	PayloadLength uint32
	/**
	 * 1byte.
	 * One byte field to represent the message type. A range of type IDs
	 * (1-7) are reserved for protocol control messages.
	 */
	Message_type uint8
	/**
	 * 4bytes.
	 * Four-byte field that identifies the stream of the message. These
	 * bytes are set in big-endian format.
	 */
	Stream_id uint32

	/**
	 * Four-byte field that contains a Timestamp of the message.
	 * The 4 bytes are packed in the big-endian order.
	 * @remark, used as calc Timestamp when decode and encode time.
	 * @remark, we use 64bits for large time for jitter detect and hls.
	 */
	Timestamp uint64

	/**
	 * get the perfered cid(chunk stream id) which sendout over.
	 * set at decoding, and canbe used for directly send message,
	 * for example, dispatch to all connections.
	 */
	Perfer_csid uint32
}

func (h *MessageHeader) Is_audio() bool {
	return RTMP_MSG_AudioMessage == h.Message_type
}
func (h *MessageHeader) Is_video() bool {
	return RTMP_MSG_VideoMessage == h.Message_type

}
func (h *MessageHeader) IsAmf0Command() bool {
	return RTMP_MSG_AMF0CommandMessage == h.Message_type

}
func (h *MessageHeader) IsAmf0Data() bool {
	return RTMP_MSG_AMF0DataMessage == h.Message_type

}
func (h *MessageHeader) IsAmf3Command() bool {
	return RTMP_MSG_AMF3CommandMessage == h.Message_type

}
func (h *MessageHeader) IsAmf3Data() bool {
	return RTMP_MSG_AMF3DataMessage == h.Message_type

}
func (h *MessageHeader) IsWindowAckledgementSize() bool {
	return RTMP_MSG_WindowAcknowledgementSize == h.Message_type

}
func (h *MessageHeader) IsAckledgement() bool {
	return RTMP_MSG_Acknowledgement == h.Message_type

}
func (h *MessageHeader) IsSetChunkSize() bool {
	return RTMP_MSG_SetChunkSize == h.Message_type

}
func (h *MessageHeader) IsUserControlMessage() bool {
	return RTMP_MSG_UserControlMessage == h.Message_type

}
func (h *MessageHeader) IsSetPeerBandwidth() bool {
	return RTMP_MSG_SetPeerBandwidth == h.Message_type

}
func (h *MessageHeader) IsAggregate() bool {
	return RTMP_MSG_AggregateMessage == h.Message_type
}

/**
 * create a amf0 script header, set the size and stream_id.
 */
func (h *MessageHeader) InitializeAmf0Script(payload_len uint32, stream_id uint32) {
	h.Message_type = RTMP_MSG_AMF0DataMessage
	h.PayloadLength = payload_len
	h.Timestamp_delta = 0
	h.Timestamp = 0
	h.Stream_id = stream_id

	// amf0 script use connection2 chunk-id
	h.Perfer_csid = RTMP_CID_OverConnection2
}

/**
 * create a audio header, set the size, timestamp and stream_id.
 */
func (h *MessageHeader) InitializeAudio(payload_size uint32, time uint32, stream_id uint32) {
	h.Message_type = RTMP_MSG_AudioMessage
	h.PayloadLength = payload_size
	h.Timestamp_delta = time
	h.Timestamp = uint64(time)
	h.Stream_id = stream_id

	// audio chunk-id
	h.Perfer_csid = RTMP_CID_Audio
}

/**
 * create a video header, set the size, timestamp and stream_id.
 */
func (h *MessageHeader) InitializeVideo(payloadSize uint32, time uint32, streamId uint32) {
	h.Message_type = RTMP_MSG_VideoMessage
	h.PayloadLength = payloadSize
	h.Timestamp_delta = time
	h.Timestamp = uint64(time)
	h.Stream_id = streamId

	// video chunk-id
	h.Perfer_csid = RTMP_CID_Video
}

/* message is raw data RTMP message, bytes oriented*/
type Message struct {
	Header MessageHeader
	/**
	 * current message parsed SizeTmp,
	 *       SizeTmp <= header.payload_length
	 * for the payload maybe sent in multiple chunks.
	 * when finish recv whole msg, it will be reset to 0.
	 */
	SizeTmp uint32
	/**
	 * the Payload of message, can not know about the detail of Payload,
	 * user must use decode_message to get concrete packet.
	 * @remark, not all message Payload can be decoded to packet. for example,
	 *       video/audio packet use raw bytes, no video/audio packet.
	 */
	Payload []uint8
}
