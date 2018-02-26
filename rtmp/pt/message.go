package pt

// MessageHeader message header.
type MessageHeader struct {
	// 3bytes.
	// Three-byte field that contains a timestamp delta of the message.
	// The 4 bytes are packed in the big-endian order.
	// only used for decoding message from chunk stream.
	TimestampDelta uint32

	// 3bytes.
	// Three-byte field that represents the size of the payload in bytes.
	// It is set in big-endian format.
	PayloadLength uint32

	// 1byte.
	// One byte field to represent the message type. A range of type IDs
	// (1-7) are reserved for protocol control messages.
	MessageType uint8

	// 4bytes.
	// Four-byte field that identifies the stream of the message. These
	// bytes are set in big-endian format.
	StreamID uint32

	// Four-byte field that contains a Timestamp of the message.
	// The 4 bytes are packed in the big-endian order.
	// @remark, used as calc Timestamp when decode and encode time.
	// @remark, we use 64bits for large time for jitter detect and hls.
	Timestamp uint64

	// get the perfered cid(chunk stream id) which sendout over.
	// set at decoding, and canbe used for directly send message,
	// for example, dispatch to all connections.
	PerferCsid uint32
}

// IsAudio .
func (h *MessageHeader) IsAudio() bool {
	return RtmpMsgAudioMessage == h.MessageType
}

// IsVideo .
func (h *MessageHeader) IsVideo() bool {
	return RtmpMsgVideoMessage == h.MessageType
}

// IsAmf0Command .
func (h *MessageHeader) IsAmf0Command() bool {
	return RtmpMsgAmf0CommandMessage == h.MessageType
}

// IsAmf0Data .
func (h *MessageHeader) IsAmf0Data() bool {
	return RtmpMsgAmf0DataMessage == h.MessageType
}

// IsAmf3Command .
func (h *MessageHeader) IsAmf3Command() bool {
	return RtmpMsgAmf3CommandMessage == h.MessageType
}

// IsAmf3Data .
func (h *MessageHeader) IsAmf3Data() bool {
	return RtmpMsgAmf3DataMessage == h.MessageType
}

// IsWindowAckledgementSize .
func (h *MessageHeader) IsWindowAckledgementSize() bool {
	return RtmpMsgWindowAcknowledgementSize == h.MessageType
}

// IsAckledgement .
func (h *MessageHeader) IsAckledgement() bool {
	return RtmpMsgAcknowledgement == h.MessageType
}

// IsSetChunkSize .
func (h *MessageHeader) IsSetChunkSize() bool {
	return RtmpMsgSetChunkSize == h.MessageType
}

// IsUserControlMessage .
func (h *MessageHeader) IsUserControlMessage() bool {
	return RtmpMsgUserControlMessage == h.MessageType
}

// IsSetPeerBandwidth .
func (h *MessageHeader) IsSetPeerBandwidth() bool {
	return RtmpMsgSetPeerBandwidth == h.MessageType
}

// IsAggregate .
func (h *MessageHeader) IsAggregate() bool {
	return RtmpMsgAggregateMessage == h.MessageType
}

// InitializeAmf0Script create a amf0 script header, set the size and stream_id.
func (h *MessageHeader) InitializeAmf0Script(payloadLen uint32, streamID uint32) {
	h.MessageType = RtmpMsgAmf0DataMessage
	h.PayloadLength = payloadLen
	h.TimestampDelta = 0
	h.Timestamp = 0
	h.StreamID = streamID

	// amf0 script use connection2 chunk-id
	h.PerferCsid = RtmpCidOverConnection2
}

// InitializeAudio create a audio header, set the size, timestamp and stream_id.
func (h *MessageHeader) InitializeAudio(payloadSize uint32, time uint32, streamID uint32) {
	h.MessageType = RtmpMsgAudioMessage
	h.PayloadLength = payloadSize
	h.TimestampDelta = time
	h.Timestamp = uint64(time)
	h.StreamID = streamID

	// audio chunk-id
	h.PerferCsid = RtmpCidAudio
}

// InitializeVideo create a video header, set the size, timestamp and stream_id.
func (h *MessageHeader) InitializeVideo(payloadSize uint32, time uint32, streamID uint32) {
	h.MessageType = RtmpMsgVideoMessage
	h.PayloadLength = payloadSize
	h.TimestampDelta = time
	h.Timestamp = uint64(time)
	h.StreamID = streamID

	// video chunk-id
	h.PerferCsid = RtmpCidVideo
}

// MessagePayload message payload
type MessagePayload struct {
	// current message parsed SizeTmp,
	// SizeTmp <= header.payload_length
	// for the payload maybe sent in multiple chunks.
	// when finish recv whole msg, it will be reset to 0.
	SizeTmp uint32

	// the Payload of message, can not know about the detail of Payload,
	// user must use decode_message to get concrete packet.
	// not all message Payload can be decoded to packet. for example,
	// video/audio packet use raw bytes, no video/audio packet.
	Payload []uint8
}

// Message message is raw data RTMP message, bytes oriented
type Message struct {
	Header  MessageHeader
	Payload MessagePayload
}
