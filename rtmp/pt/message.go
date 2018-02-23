package pt

//message header.
type MessageHeader struct {
	/**
	 * 3bytes.
	 * Three-byte field that contains a timestamp delta of the message.
	 * The 4 bytes are packed in the big-endian order.
	 * @remark, only used for decoding message from chunk stream.
	 */
	TimestampDelta uint32
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
	MessageType uint8
	/**
	 * 4bytes.
	 * Four-byte field that identifies the stream of the message. These
	 * bytes are set in big-endian format.
	 */
	StreamId uint32

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
	PerferCsid uint32
}

func (h *MessageHeader) IsAudio() bool {
	return RTMP_MSG_AudioMessage == h.MessageType
}
func (h *MessageHeader) IsVideo() bool {
	return RTMP_MSG_VideoMessage == h.MessageType
}

func (h *MessageHeader) IsAmf0Command() bool {
	return RTMP_MSG_AMF0CommandMessage == h.MessageType

}
func (h *MessageHeader) IsAmf0Data() bool {
	return RTMP_MSG_AMF0DataMessage == h.MessageType

}
func (h *MessageHeader) IsAmf3Command() bool {
	return RTMP_MSG_AMF3CommandMessage == h.MessageType

}
func (h *MessageHeader) IsAmf3Data() bool {
	return RTMP_MSG_AMF3DataMessage == h.MessageType

}
func (h *MessageHeader) IsWindowAckledgementSize() bool {
	return RTMP_MSG_WindowAcknowledgementSize == h.MessageType

}
func (h *MessageHeader) IsAckledgement() bool {
	return RTMP_MSG_Acknowledgement == h.MessageType

}
func (h *MessageHeader) IsSetChunkSize() bool {
	return RTMP_MSG_SetChunkSize == h.MessageType

}
func (h *MessageHeader) IsUserControlMessage() bool {
	return RTMP_MSG_UserControlMessage == h.MessageType

}
func (h *MessageHeader) IsSetPeerBandwidth() bool {
	return RTMP_MSG_SetPeerBandwidth == h.MessageType

}
func (h *MessageHeader) IsAggregate() bool {
	return RTMP_MSG_AggregateMessage == h.MessageType
}

/**
 * create a amf0 script header, set the size and stream_id.
 */
func (h *MessageHeader) InitializeAmf0Script(payload_len uint32, stream_id uint32) {
	h.MessageType = RTMP_MSG_AMF0DataMessage
	h.PayloadLength = payload_len
	h.TimestampDelta = 0
	h.Timestamp = 0
	h.StreamId = stream_id

	// amf0 script use connection2 chunk-id
	h.PerferCsid = RTMP_CID_OverConnection2
}

/**
 * create a audio header, set the size, timestamp and stream_id.
 */
func (h *MessageHeader) InitializeAudio(payload_size uint32, time uint32, stream_id uint32) {
	h.MessageType = RTMP_MSG_AudioMessage
	h.PayloadLength = payload_size
	h.TimestampDelta = time
	h.Timestamp = uint64(time)
	h.StreamId = stream_id

	// audio chunk-id
	h.PerferCsid = RTMP_CID_Audio
}

/**
 * create a video header, set the size, timestamp and stream_id.
 */
func (h *MessageHeader) InitializeVideo(payloadSize uint32, time uint32, streamId uint32) {
	h.MessageType = RTMP_MSG_VideoMessage
	h.PayloadLength = payloadSize
	h.TimestampDelta = time
	h.Timestamp = uint64(time)
	h.StreamId = streamId

	// video chunk-id
	h.PerferCsid = RTMP_CID_Video
}

type MessagePayload struct {
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

/* message is raw data RTMP message, bytes oriented*/
type Message struct {
	Header  MessageHeader
	Payload MessagePayload
}
