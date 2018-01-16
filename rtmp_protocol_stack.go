package main

const (
	/**
	 * 6.1.2. Chunk Message Header
	 * There are four different formats for the chunk message header,
	 * selected by the "chunkFmt" field in the chunk basic header.
	 */
	// 6.1.2.1. Type 0
	// Chunks of Type 0 are 11 bytes long. This type MUST be used at the
	// start of a chunk stream, and whenever the stream timestampDelta goes
	// backward (e.g., because of a backward seek).
	RTMP_FMT_TYPE0 = 0
	// 6.1.2.2. Type 1
	// Chunks of Type 1 are 7 bytes long. The message stream ID is not
	// included; this chunk takes the same stream ID as the preceding chunk.
	// Streams with variable-sized messages (for example, many video
	// formats) SHOULD use this format for the first chunk of each new
	// message after the first.
	RTMP_FMT_TYPE1 = 1
	// 6.1.2.3. Type 2
	// Chunks of Type 2 are 3 bytes long. Neither the stream ID nor the
	// message length is included; this chunk has the same stream ID and
	// message length as the preceding chunk. Streams with constant-sized
	// messages (for example, some audio and data formats) SHOULD use this
	// format for the first chunk of each message after the first.
	RTMP_FMT_TYPE2 = 2
	// 6.1.2.4. Type 3
	// Chunks of Type 3 have no header. Stream ID, message length and
	// timestampDelta delta are not present; chunks of this type take values from
	// the preceding chunk. When a single message is split into chunks, all
	// chunks of a message except the first one, SHOULD use this type. Refer
	// to example 2 in section 6.2.2. Stream consisting of messages of
	// exactly the same size, stream ID and spacing in time SHOULD use this
	// type for all chunks after chunk of Type 2. Refer to example 1 in
	// section 6.2.1. If the delta between the first message and the second
	// message is same as the time stamp of first message, then chunk of
	// type 3 would immediately follow the chunk of type 0 as there is no
	// need for a chunk of type 2 to register the delta. If Type 3 chunk
	// follows a Type 0 chunk, then timestampDelta delta for this Type 3 chunk is
	// the same as the timestampDelta of Type 0 chunk.
	RTMP_FMT_TYPE3 = 3
)

const (
	/**
	 * the chunk stream id used for some under-layer message,
	 * for example, the PC(protocol control) message.
	 */
	RTMP_CID_ProtocolControl = 0x02
	/**
	 * the AMF0/AMF3 command message, invoke method and return the result, over NetConnection.
	 * generally use 0x03.
	 */
	RTMP_CID_OverConnection = 0x03
	/**
	 * the AMF0/AMF3 command message, invoke method and return the result, over NetConnection,
	 * the midst state(we guess).
	 * rarely used, e.g. onStatus(NetStream.Play.Reset).
	 */
	RTMP_CID_OverConnection2 = 0x04
	/**
	 * the stream message(amf0/amf3), over NetStream.
	 * generally use 0x05.
	 */
	RTMP_CID_OverStream = 0x05
	/**
	 * the stream message(amf0/amf3), over NetStream, the midst state(we guess).
	 * rarely used, e.g. play("mp4:mystram.f4v")
	 */
	RTMP_CID_OverStream2 = 0x08
	/**
	 * the stream message(video), over NetStream
	 * generally use 0x06.
	 */
	RTMP_CID_Video = 0x06
	/**
	 * the stream message(audio), over NetStream.
	 * generally use 0x07.
	 */
	RTMP_CID_Audio = 0x07
)

const (
	RTMP_EXTENDED_TIMESTAMP = 0xFFFFFF
)

const (
	RTMP_DEFAULT_CHUNK_SIZE = 128
)
