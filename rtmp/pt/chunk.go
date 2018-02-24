package pt

// ChunkStream incoming chunk stream maybe interlaced,
// use the chunk stream to cache the input RTMP chunk streams.
type ChunkStream struct {
	/**
	 * represents the basic header fmt,
	 * which used to identify the variant message header type.
	 */
	Fmt uint8
	/**
	 * represents the basic header cs_id,
	 * which is the chunk stream id.
	 */
	CsID uint32
	/**
	 * cached message header
	 */
	MsgHeader MessageHeader
	/**
	 * whether the chunk message header has extended timestamp.
	 */
	HasExtendedTimestamp bool
	/**
	 * decoded msg count, to identify whether the chunk stream is fresh.
	 */
	MsgCount uint64
}
