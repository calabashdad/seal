package main

import "encoding/binary"

func (rtmp *RtmpConn) SendMsg(msg *MessageStream) (err error) {

	var payloadOffset uint32

	for {

		if payloadOffset >= msg.header.length {
			break
		}

		var headerOffset uint32
		var header [16]uint8 // 1(max) basic + 11(max) msg header + 4(max) extend timestamp

		if 0 == payloadOffset {
			// new chunk. fmt type 0. msg header 11 bytes.

			//fmt
			fmt := 0x00 | (msg.header.preferCsId & 0x3f)
			header[headerOffset] = uint8(fmt)
			headerOffset += 1

			//message header -- timestamp
			if msg.header.timestamp < RTMP_EXTENDED_TIMESTAMP {
				// big endian
				header[headerOffset] = uint8(uint32(msg.header.timestamp) << 8 >> 24)
				headerOffset += 1
				header[headerOffset] = uint8(uint32(msg.header.timestamp) << 16 >> 24)
				headerOffset += 1
				header[headerOffset] = uint8(uint32(msg.header.timestamp) << 24 >> 24)
				headerOffset += 1
			} else {
				header[headerOffset] = 0xff
				headerOffset += 1
				header[headerOffset] = 0xff
				headerOffset += 1
				header[headerOffset] = 0xff
				headerOffset += 1
			}

			//message -- length. big endian
			header[headerOffset] = uint8(msg.header.length << 8 >> 24)
			headerOffset += 1
			header[headerOffset] = uint8(msg.header.length << 16 >> 24)
			headerOffset += 1
			header[headerOffset] = uint8(msg.header.length << 24 >> 24)
			headerOffset += 1

			//message -- message type
			header[headerOffset] = msg.header.typeId
			headerOffset += 1

			//message -- stream id. little endian
			binary.LittleEndian.PutUint32(header[headerOffset:headerOffset+4], msg.header.streamId)
			headerOffset += 4

			//extend timestamp
			if msg.header.timestamp > RTMP_EXTENDED_TIMESTAMP {
				binary.BigEndian.PutUint32(header[headerOffset:headerOffset+4], uint32(msg.header.timestamp))
				headerOffset += 4
			}

			//
		} else {

			// fmt type 3. no msg header.

			//fmt
			fmt := 0xc0 | (msg.header.preferCsId & 0x3f)
			header[headerOffset] = uint8(fmt)
			headerOffset += 1

			//extend timestamp
			if msg.header.timestamp > RTMP_EXTENDED_TIMESTAMP {
				binary.BigEndian.PutUint32(header[headerOffset:headerOffset+4], uint32(msg.header.timestamp))
				headerOffset += 4
			}
		}

		var payloadSizeSendThisTime uint32
		if msg.header.length-payloadOffset > rtmp.chunkSize {
			payloadSizeSendThisTime = rtmp.chunkSize
		} else {
			payloadSizeSendThisTime = msg.header.length - payloadOffset
		}

		//send header
		err = rtmp.SendBytes(header[:headerOffset])
		if err != nil {
			break
		}

		//send payload
		err = rtmp.SendBytes(msg.payload[payloadOffset : payloadOffset+payloadSizeSendThisTime])
		if err != nil {
			break
		}

		payloadOffset += payloadSizeSendThisTime
	}

	if err != nil {
		return
	}

	return
}
