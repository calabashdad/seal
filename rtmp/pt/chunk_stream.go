package pt

/**
* incoming chunk stream maybe interlaced,
* use the chunk stream to cache the input RTMP chunk streams.
 */
type ChunkStream struct {
	/**
	 * represents the basic header fmt,
	 * which used to identify the variant message header type.
	 */
	Header_fmt uint8
	/**
	 * represents the basic header cs_id,
	 * which is the chunk stream id.
	 */
	Cs_id uint32
	/**
	 * cached message header
	 */
	Msg_header MessageHeader
	/**
	 * whether the chunk message header has extended timestamp.
	 */
	Extended_timestamp bool
	/**
	 * partially read message.
	 */
	Msg *Message
	/**
	 * decoded msg count, to identify whether the chunk stream is fresh.
	 */
	Msg_count uint64
}

func (chunk *ChunkStream) GotEntireMsg() bool {
	return (chunk.Msg.Header.Payload_length == chunk.Msg.Size)
}
