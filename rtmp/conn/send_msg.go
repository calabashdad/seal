package conn

import (
	"UtilsTools/identify_panic"
	"encoding/binary"
	"log"
	"seal/rtmp/protocol"
)

func (rc *RtmpConn) SendMsg(msg *protocol.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	if nil == msg {
		return
	}

	// ensure the basic header is 1bytes. make simple.
	if msg.Header.Perfer_csid < 2 {
		msg.Header.Perfer_csid = protocol.RTMP_CID_ProtocolControl
	}

	//current position of payload send.
	var payload_offset uint32

	// always write the header event payload is empty.
	for {
		if payload_offset >= msg.Header.Payload_length {
			break
		}

		var header_offset uint32
		var header [RTMP_MAX_FMT0_HEADER_SIZE]uint8

		if 0 == payload_offset {
			// write new chunk stream header, fmt is 0
			header[header_offset] = 0x00 | uint8(msg.Header.Perfer_csid&0x3f)
			header_offset += 1

			// chunk message header, 11 bytes
			// timestamp, 3bytes, big-endian
			timestamp := msg.Header.Timestamp
			if timestamp < protocol.RTMP_EXTENDED_TIMESTAMP {
				header[header_offset] = uint8((timestamp & 0x00ff0000) >> 16)
				header_offset += 1
				header[header_offset] = uint8((timestamp & 0x0000ff00) >> 8)
				header_offset += 1
				header[header_offset] = uint8(timestamp & 0x000000ff)
				header_offset += 1
			} else {
				header[header_offset] = 0xff
				header_offset += 1
				header[header_offset] = 0xff
				header_offset += 1
				header[header_offset] = 0xff
				header_offset += 1
			}

			// message_length, 3bytes, big-endian
			payload_lengh := msg.Header.Payload_length
			header[header_offset] = uint8((payload_lengh & 0x00ff0000) >> 16)
			header_offset += 1
			header[header_offset] = uint8((payload_lengh & 0x0000ff00) >> 8)
			header_offset += 1
			header[header_offset] = uint8((payload_lengh & 0x000000ff))
			header_offset += 1

			// message_type, 1bytes
			header[header_offset] = msg.Header.Message_type
			header_offset += 1

			// stream id, 4 bytes, little-endian
			binary.LittleEndian.PutUint32(header[header_offset:header_offset+4], msg.Header.Stream_id)
			header_offset += 4

			// chunk extended timestamp header, 0 or 4 bytes, big-endian
			if timestamp >= protocol.RTMP_EXTENDED_TIMESTAMP {
				binary.BigEndian.PutUint32(header[header_offset:header_offset+4], uint32(timestamp))
				header_offset += 4
			}

		} else {
			// write no message header chunk stream, fmt is 3
			// @remark, if perfer_cid > 0x3F, that is, use 2B/3B chunk header,
			// rollback to 1B chunk header.

			// fmt is 3
			header[header_offset] = 0xc0 | uint8(msg.Header.Perfer_csid&0x3f)
			header_offset += 1

			// chunk extended timestamp header, 0 or 4 bytes, big-endian
			// 6.1.3. Extended Timestamp
			// This field is transmitted only when the normal time stamp in the
			// chunk message header is set to 0x00ffffff. If normal time stamp is
			// set to any value less than 0x00ffffff, this field MUST NOT be
			// present. This field MUST NOT be present if the timestamp field is not
			// present. Type 3 chunks MUST NOT have this field.
			// adobe changed for Type3 chunk:
			//        FMLE always sendout the extended-timestamp,
			//        must send the extended-timestamp to FMS,
			//        must send the extended-timestamp to flash-player.
			timestamp := msg.Header.Timestamp
			if timestamp >= protocol.RTMP_EXTENDED_TIMESTAMP {
				binary.BigEndian.PutUint32(header[header_offset:header_offset+4], uint32(timestamp))
				header_offset += 4
			}
		}

		//send header
		err = rc.TcpConn.SendBytes(header[:header_offset])
		if err != nil {
			log.Println("send msg header failed.")
			return
		}

		//payload
		payload_size := msg.Header.Payload_length - payload_offset
		if payload_size > rc.Out_chunk_size {
			payload_size = rc.Out_chunk_size
		}

		err = rc.TcpConn.SendBytes(msg.Payload[payload_offset : payload_offset+payload_size])
		if err != nil {
			log.Println("send msg payload failed.")
			return
		}

		payload_offset += payload_size
	}

	return
}
