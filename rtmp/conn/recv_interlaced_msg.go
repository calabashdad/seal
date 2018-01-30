package conn

import (
	"UtilsTools/identify_panic"
	"encoding/binary"
	"fmt"
	"log"
	"seal/rtmp/pt"
)

func (rc *RtmpConn) RecvInterlacedMsg(msg **pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	// chunk stream basic header.
	var chunk_fmt uint8
	var cs_id uint32
	var chunk_header_size uint32

	err = rc.ReadBasicHeader(&chunk_fmt, &cs_id, &chunk_header_size)
	if err != nil {
		return
	}

	chunk := rc.ChunkStreams[cs_id]
	if nil == chunk {
		chunk = &pt.ChunkStream{
			Cs_id: cs_id,
		}

		rc.ChunkStreams[cs_id] = chunk
	}

	//read msg header
	var msg_header_size uint32
	err = rc.ReadMessageHeader(chunk, chunk_fmt, chunk_header_size, &msg_header_size)
	if err != nil {
		return
	}

	//read msg payload
	err = rc.ReadMsgPayload(chunk)
	if err != nil {
		return
	}

	if chunk.GotEntireMsg() {
		*msg = chunk.Msg
		//may be the msg cache can be reused by multi chunks, so when recv an entire msg, reset to zero.
		(*msg).Size = 0
	}

	return
}

/**
* 6.1.1. Chunk Basic Header
* The Chunk Basic Header encodes the chunk stream ID and the chunk
* type(represented by fmt field in the figure below). Chunk type
* determines the format of the encoded message header. Chunk Basic
* Header field may be 1, 2, or 3 bytes, depending on the chunk stream
* ID.
*
* The bits 0–5 (least significant) in the chunk basic header represent
* the chunk stream ID.
*
* Chunk stream IDs 2-63 can be encoded in the 1-byte version of this
* field.
*    0 1 2 3 4 5 6 7
*   +-+-+-+-+-+-+-+-+
*   |fmt|   cs id   |
*   +-+-+-+-+-+-+-+-+
*   Figure 6 Chunk basic header 1
*
* Chunk stream IDs 64-319 can be encoded in the 2-byte version of this
* field. ID is computed as (the second byte + 64).
*   0                   1
*   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5
*   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
*   |fmt|    0      | cs id - 64    |
*   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
*   Figure 7 Chunk basic header 2
*
* Chunk stream IDs 64-65599 can be encoded in the 3-byte version of
* this field. ID is computed as ((the third byte)*256 + the second byte
* + 64).
*    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3
*   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
*   |fmt|     1     |         cs id - 64            |
*   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
*   Figure 8 Chunk basic header 3
*
* cs id: 6 bits
* fmt: 2 bits
* cs id - 64: 8 or 16 bits
*
* Chunk stream IDs with values 64-319 could be represented by both 2-
* byte version and 3-byte version of this field.
 */
func (rc *RtmpConn) ReadBasicHeader(header_fmt *uint8, cs_id *uint32, chunk_header_size *uint32) (err error) {

	var buf [3]uint8
	err = rc.TcpConn.ExpectBytesFull(buf[:1], 1)
	if err != nil {
		return
	}

	*header_fmt = (buf[0] >> 6) & 0x03
	*cs_id = uint32(buf[0] & 0x3f)
	*chunk_header_size = 1

	// 2-63, 1B chunk header
	if *cs_id > 1 {
		return
	}

	// 64-319, 2B chunk header
	if 0 == *cs_id {
		err = rc.TcpConn.ExpectBytesFull(buf[1:], 1)
		if err != nil {
			return
		}

		*cs_id = 64
		*cs_id += uint32(buf[1])
		*chunk_header_size = 2
	} else if 1 == *cs_id {
		// 64-65599, 3B chunk header
		err = rc.TcpConn.ExpectBytesFull(buf[1:], 2)
		if err != nil {
			return
		}

		*cs_id = 64
		*cs_id += uint32(buf[1])
		*cs_id += uint32(buf[2]) * 256
		*chunk_header_size = 3
	} else {
		err = fmt.Errorf("invalid cs id.")
		return
	}

	return
}

/**
* parse the message header.
*   3bytes: timestamp delta,    fmt=0,1,2
*   3bytes: payload length,     fmt=0,1
*   1bytes: message type,       fmt=0,1
*   4bytes: stream id,          fmt=0
* where:
*   fmt=0, 0x0X
*   fmt=1, 0x4X
*   fmt=2, 0x8X
*   fmt=3, 0xCX
 */
func (rc *RtmpConn) ReadMessageHeader(chunk *pt.ChunkStream, chunkFmt uint8, chunkHeaderSize uint32, msgHeaderSize *uint32) (err error) {

	/**
	 * we should not assert anything about fmt, for the first packet.
	 * (when first packet, the chunk->msg is NULL).
	 * the fmt maybe 0/1/2/3, the FMLE will send a 0xC4 for some audio packet.
	 * the previous packet is:
	 *     04                // fmt=0, cid=4
	 *     00 00 1a          // timestamp=26
	 *     00 00 9d          // payload_length=157
	 *     08                // message_type=8(audio)
	 *     01 00 00 00       // stream_id=1
	 * the current packet maybe:
	 *     c4             // fmt=3, cid=4
	 * it's ok, for the packet is audio, and timestamp delta is 26.
	 * the current packet must be parsed as:
	 *     fmt=0, cid=4
	 *     timestamp=26+26=52
	 *     payload_length=157
	 *     message_type=8(audio)
	 *     stream_id=1
	 * so we must update the timestamp even fmt=3 for first packet.
	 */
	// fresh packet used to update the timestamp even fmt=3 for first packet.
	// fresh packet always means the chunk is the first one of message.

	var is_first_chunk_of_msg bool
	if nil == chunk.Msg {
		is_first_chunk_of_msg = true
	} else {
		if 0 == chunk.Msg.Size {
			is_first_chunk_of_msg = true
		} else {
			is_first_chunk_of_msg = false
		}
	}

	//when a chunk stream is fresh, the fmt must be 0, a new stream.
	if 0 == chunk.Msg_count && chunkFmt != pt.RTMP_FMT_TYPE0 {
		// for librtmp, if ping, it will send a fresh stream with fmt=1,
		// 0x42             where: fmt=1, cid=2, pt contorl user-control message
		// 0x00 0x00 0x00   where: timestamp=0
		// 0x00 0x00 0x06   where: payload_length=6
		// 0x04             where: message_type=4(pt control user-control message)
		// 0x00 0x06            where: event Ping(0x06)
		// 0x00 0x00 0x0d 0x0f  where: event data 4bytes ping timestamp.
		if pt.RTMP_CID_ProtocolControl == chunk.Cs_id && pt.RTMP_FMT_TYPE1 == chunkFmt {
			log.Println("accept cid=2, fmt=1 to make librtmp happy.")
		} else {
			// must be a RTMP pt level error.
			err = fmt.Errorf("chunk stream is fresh, fmt must be RTMP_FMT_TYPE0, actual is %d", chunkFmt)
			return
		}
	}

	// when exists cache msg, means got an partial message,
	// the fmt must not be type0 which means new message.
	if nil != chunk.Msg && chunk.Msg.Size > 0 && pt.RTMP_FMT_TYPE0 == chunkFmt {
		err = fmt.Errorf("chunk stream exists, fmt must not be RTMP_FMT_TYPE0, actual is %d", chunkFmt)
		return
	}

	if nil == chunk.Msg {
		chunk.Msg = &pt.Message{}
	}

	switch chunkFmt {
	case 0:
		*msgHeaderSize = 11
	case 1:
		*msgHeaderSize = 7
	case 2:
		*msgHeaderSize = 3
	case 3:
		*msgHeaderSize = 0
	default:
		err = fmt.Errorf("invalid chunk fmt when calc msg header size.")
	}

	var msg_header_buf [12]uint8 //max is 11

	err = rc.TcpConn.ExpectBytesFull(msg_header_buf[:], *msgHeaderSize)
	if err != nil {
		return
	}

	/**
	 * parse the message header.
	 *   3bytes: timestamp delta,    fmt=0,1,2
	 *   3bytes: payload length,     fmt=0,1
	 *   1bytes: message type,       fmt=0,1
	 *   4bytes: stream id,          fmt=0
	 * where:
	 *   fmt=0, 0x0X
	 *   fmt=1, 0x4X
	 *   fmt=2, 0x8X
	 *   fmt=3, 0xCX
	 */

	if chunkFmt <= pt.RTMP_FMT_TYPE2 {
		chunk.Msg_header.Timestamp_delta = 0

		var offset uint32
		chunk.Msg_header.Timestamp_delta |= (uint32(msg_header_buf[offset]) << 16)
		offset++
		chunk.Msg_header.Timestamp_delta |= (uint32(msg_header_buf[offset]) << 8)
		offset++
		chunk.Msg_header.Timestamp_delta |= (uint32(msg_header_buf[offset]))
		offset++

		// fmt: 0
		// timestamp: 3 bytes
		// If the timestamp is greater than or equal to 16777215
		// (hexadecimal 0x00ffffff), this value MUST be 16777215, and the
		// ‘extended timestamp header’ MUST be present. Otherwise, this value
		// SHOULD be the entire timestamp.
		//
		// fmt: 1 or 2
		// timestamp delta: 3 bytes
		// If the delta is greater than or equal to 16777215 (hexadecimal
		// 0x00ffffff), this value MUST be 16777215, and the ‘extended
		// timestamp header’ MUST be present. Otherwise, this value SHOULD be
		// the entire delta.

		chunk.Extended_timestamp = (chunk.Msg_header.Timestamp_delta >= pt.RTMP_EXTENDED_TIMESTAMP)
		if !chunk.Extended_timestamp {
			// Extended timestamp: 0 or 4 bytes
			// This field MUST be sent when the normal timsestamp is set to
			// 0xffffff, it MUST NOT be sent if the normal timestamp is set to
			// anything else. So for values less than 0xffffff the normal
			// timestamp field SHOULD be used in which case the extended timestamp
			// MUST NOT be present. For values greater than or equal to 0xffffff
			// the normal timestamp field MUST NOT be used and MUST be set to
			// 0xffffff and the extended timestamp MUST be sent.

			if pt.RTMP_FMT_TYPE0 == chunkFmt {
				// 6.1.2.1. Type 0
				// For a type-0 chunk, the absolute timestamp of the message is sent
				// here.
				chunk.Msg_header.Timestamp = uint64(chunk.Msg_header.Timestamp_delta)
			} else {
				// 6.1.2.2. Type 1
				// 6.1.2.3. Type 2
				// For a type-1 or type-2 chunk, the difference between the previous
				// chunk's timestamp and the current chunk's timestamp is sent here.
				chunk.Msg_header.Timestamp += uint64(chunk.Msg_header.Timestamp_delta)
			}
		}

		if chunkFmt <= pt.RTMP_FMT_TYPE1 {
			var payloadLength uint32

			payloadLength |= (uint32(msg_header_buf[offset]) << 16)
			offset++
			payloadLength |= (uint32(msg_header_buf[offset]) << 8)
			offset++
			payloadLength |= (uint32(msg_header_buf[offset]))
			offset++

			// for a message, if msg exists in cache, the size must not changed.
			// always use the actual msg size to compare, for the cache payload length can changed,
			// for the fmt type1(stream_id not changed), user can change the payload
			// length(it's not allowed in the continue chunks).
			if !is_first_chunk_of_msg && chunk.Msg_header.Payload_length != payloadLength {
				err = fmt.Errorf("msg exists in chunk cache, payload size=%d is not equal to msg header, should be=%d",
					payloadLength, chunk.Msg_header.Payload_length)
				return
			}

			chunk.Msg_header.Payload_length = payloadLength
			chunk.Msg_header.Message_type = msg_header_buf[offset]
			offset += 1

			if pt.RTMP_FMT_TYPE0 == chunkFmt {
				chunk.Msg_header.Stream_id = binary.LittleEndian.Uint32(msg_header_buf[offset : offset+4])
				offset += 4

			} else {
				//header read completed
			}
		} else {
			//header read completed
		}
	} else {
		// update the timestamp even fmt=3 for first chunk packet
		if is_first_chunk_of_msg && !chunk.Extended_timestamp {
			chunk.Msg_header.Timestamp += uint64(chunk.Msg_header.Timestamp_delta)
		}
	}

	// read extended-timestamp
	if chunk.Extended_timestamp {
		var extendTimestampBuf [4]uint8
		err = rc.TcpConn.ExpectBytesFull(extendTimestampBuf[:], 4)
		if err != nil {
			return
		}

		extendTimestamp := binary.BigEndian.Uint32(extendTimestampBuf[0:4])

		// always use 31bits timestamp, for some server may use 32bits extended timestamp.
		extendTimestamp &= 0x7fffffff

		/**
		* RTMP specification and ffmpeg/librtmp is false,
		* but, adobe changed the specification, so flash/FMLE/FMS always true.
		* default to true to support flash/FMLE/FMS.
		*
		* ffmpeg/librtmp may donot send this filed, need to detect the value.
		* compare to the chunk timestamp, which is set by chunk message header
		* type 0,1 or 2.
		*
		* @remark, nginx send the extended-timestamp in sequence-header,
		* and timestamp delta in continue C1 chunks, and so compatible with ffmpeg,
		* that is, there is no continue chunks and extended-timestamp in nginx-rtmp.
		*
		* @remark, seal always send the extended-timestamp, to keep simple,
		* and compatible with adobe products.
		 */
		chunk_timestamp := chunk.Msg_header.Timestamp

		/**
		* if chunk_timestamp<=0, the chunk previous packet has no extended-timestamp,
		* always use the extended timestamp.
		 */
		/**
		* about the is_first_chunk_of_msg.
		* @remark, for the first chunk of message, always use the extended timestamp.
		 */
		if !is_first_chunk_of_msg && chunk_timestamp > 0 && chunk_timestamp != uint64(extendTimestamp) {
			//("no 4bytes extended timestamp in the continued chunk");
		} else {
			chunk.Msg_header.Timestamp = uint64(extendTimestamp)
		}
	}

	// the extended-timestamp must be unsigned-int,
	//         24bits timestamp: 0xffffff = 16777215ms = 16777.215s = 4.66h
	//         32bits timestamp: 0xffffffff = 4294967295ms = 4294967.295s = 1193.046h = 49.71d
	// because the rtmp pt says the 32bits timestamp is about "50 days":
	//         3. Byte Order, Alignment, and Time Format
	//                Because timestamps are generally only 32 bits long, they will roll
	//                over after fewer than 50 days.
	//
	// but, its sample says the timestamp is 31bits:
	//         An application could assume, for example, that all
	//        adjacent timestamps are within 2^31 milliseconds of each other, so
	//        10000 comes after 4000000000, while 3000000000 comes before
	//        4000000000.
	// and flv specification says timestamp is 31bits:
	//        Extension of the Timestamp field to form a SI32 value. This
	//        field represents the upper 8 bits, while the previous
	//        Timestamp field represents the lower 24 bits of the time in
	//        milliseconds.
	// in a word, 31bits timestamp is ok.
	// convert extended timestamp to 31bits.
	chunk.Msg_header.Timestamp &= 0x7fffffff

	// copy header to msg
	chunk.Msg.Header = chunk.Msg_header

	// increase the msg count, the chunk stream can accept fmt=1/2/3 message now.
	chunk.Msg_count++

	return
}

func (rc *RtmpConn) ReadMsgPayload(chunk *pt.ChunkStream) (err error) {

	payloadSize := chunk.Msg_header.Payload_length - chunk.Msg.Size

	if payloadSize > rc.InChunkSize {
		payloadSize = rc.InChunkSize
	}

	if nil == chunk.Msg.Payload {
		chunk.Msg.Payload = rc.Pool.GetMem(chunk.Msg.Header.Payload_length)
		if nil == chunk.Msg.Payload {
			err = fmt.Errorf("alloc msg payload space failed.")
			return
		}
	}

	err = rc.TcpConn.ExpectBytesFull(chunk.Msg.Payload[chunk.Msg.Size:], payloadSize)
	if err != nil {
		return
	}

	chunk.Msg.Size += payloadSize

	return
}
