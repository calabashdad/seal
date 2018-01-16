package main

import (
	"UtilsTools/identify_panic"
	"encoding/binary"
	"fmt"
	"log"
)

type ChunkStruct struct {
	chunkFmt uint8
	csId     uint32
	//msg count of this chunk
	msgCount      uint64
	msgHeaderSize uint32
	msgHeader     struct {
		timestampDelta uint32
		timestamp      uint64
		msgLength      uint32
		msgTypeid      uint8
		msgStreamId    uint32
	}
	hasExtendTimestamp bool
	extendTimeStamp    uint32
	msgPayload         []uint8
	msgPayloadSize     uint32

	chunkSize uint32 //default is RTMP_DEFAULT_CHUNK_SIZE.
}

func (rtmp *RtmpSession) ExpectMsg() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}
	}()

	for {

		//read basic header
		var buf [3]uint8

		err = rtmp.ExpectBytes(1, buf[:1])
		if err != nil {
			return
		}

		chunk_fmt := buf[0] & 0xc0
		csid := uint32(buf[0] & 0x3f)

		switch csid {
		case 0:
			//csId 2 bytes. 64-319
			err = rtmp.ExpectBytes(1, buf[1:2])
			csid = uint32(64 + buf[1])
		case 1:
			//csId 3 bytes. 64-65599
			err = rtmp.ExpectBytes(2, buf[1:3])
			csid = uint32(64) + uint32(buf[1]) + uint32(buf[2])*256
		}

		if err != nil {
			return
		}

		chunk, ok := rtmp.chunks[csid]
		if !ok {
			chunk = ChunkStruct{
				chunkSize: RTMP_DEFAULT_CHUNK_SIZE,
			}
		}

		chunk.chunkFmt = chunk_fmt
		chunk.csId = csid

		rtmp.chunks[csid] = chunk

		//read message header
		if 0 == chunk.msgCount && chunk.chunkFmt != RTMP_FMT_TYPE0 {
			if RTMP_CID_ProtocolControl == chunk.csId && RTMP_FMT_TYPE1 == chunk.chunkFmt {
				// for librtmp, if ping, it will send a fresh stream with fmt=1,
				// 0x42             where: fmt=1, cid=2, protocol contorl user-control message
				// 0x00 0x00 0x00   where: timestamp=0
				// 0x00 0x00 0x06   where: payload_length=6
				// 0x04             where: message_type=4(protocol control user-control message)
				// 0x00 0x06            where: event Ping(0x06)
				// 0x00 0x00 0x0d 0x0f  where: event data 4bytes ping timestamp.
				log.Println("rtmp session, accept cid=2, chunkFmt=1 , it's a valid chunk format, for librtmp.")
			} else {
				err = fmt.Errorf("chunk start error, must be RTMP_FMT_TYPE0")
				break
			}
		}

		if chunk.msgPayloadSize > 0 && RTMP_FMT_TYPE0 == chunk.chunkFmt {
			err = fmt.Errorf("when msg count > 0, chunk fmt is not allowed to be RTMP_FMT_TYPE0.")
			break
		}

		switch chunk.chunkFmt {
		case RTMP_FMT_TYPE0:
			chunk.msgHeaderSize = 11
		case RTMP_FMT_TYPE1:
			chunk.msgHeaderSize = 7
		case RTMP_FMT_TYPE2:
			chunk.msgHeaderSize = 3
		case RTMP_FMT_TYPE3:
			chunk.msgHeaderSize = 0
		}

		var msgHeader [11]uint8
		err = rtmp.ExpectBytes(chunk.msgHeaderSize, msgHeader[:])
		if err != nil {
			break
		}

		//parse msg header
		//*   3bytes: timestamp delta,    fmt=0,1,2
		//*   3bytes: payload length,     fmt=0,1
		//*   1bytes: message type,       fmt=0,1
		//*   4bytes: stream id,          fmt=0
		switch chunk.chunkFmt {
		case RTMP_FMT_TYPE0:
			chunk.msgHeader.timestampDelta = uint32(msgHeader[0])<<16 + uint32(msgHeader[1])<<8 + uint32(msgHeader[2])
			if chunk.msgHeader.timestampDelta >= RTMP_EXTENDED_TIMESTAMP {
				chunk.hasExtendTimestamp = true
			} else {
				chunk.hasExtendTimestamp = false
				// For a type-0 chunk, the absolute timestamp of the message is sent here.
				chunk.msgHeader.timestamp = uint64(chunk.msgHeader.timestampDelta)
			}

			payloadLength := uint32(msgHeader[3])<<16 + uint32(msgHeader[4])<<8 + uint32(msgHeader[5])
			if chunk.msgPayloadSize > 0 && payloadLength != chunk.msgHeader.msgLength {
				err = fmt.Errorf("RTMP_FMT_TYPE0: msg has in chunk, msg size can not change.")
			}

			chunk.msgHeader.msgLength = payloadLength

			chunk.msgHeader.msgTypeid = msgHeader[6]

			chunk.msgHeader.msgStreamId = binary.BigEndian.Uint32(msgHeader[7:11])

		case RTMP_FMT_TYPE1:
			chunk.msgHeader.timestampDelta = uint32(msgHeader[0])<<16 + uint32(msgHeader[1])<<8 + uint32(msgHeader[2])
			if chunk.msgHeader.timestampDelta >= RTMP_EXTENDED_TIMESTAMP {
				chunk.hasExtendTimestamp = true
			} else {
				chunk.hasExtendTimestamp = false
				chunk.msgHeader.timestamp += uint64(chunk.msgHeader.timestampDelta)
			}

			payloadLength := uint32(msgHeader[3])<<16 + uint32(msgHeader[4])<<8 + uint32(msgHeader[5])
			if chunk.msgPayloadSize > 0 && payloadLength != chunk.msgHeader.msgLength {
				err = fmt.Errorf("RTMP_FMT_TYPE1: msg has in chunk, msg size can not change.")
			}

			chunk.msgHeader.msgLength = payloadLength

			chunk.msgHeader.msgTypeid = msgHeader[6]

		case RTMP_FMT_TYPE2:
			chunk.msgHeader.timestampDelta = uint32(msgHeader[0])<<16 + uint32(msgHeader[1])<<8 + uint32(msgHeader[2])
			if chunk.msgHeader.timestampDelta >= RTMP_EXTENDED_TIMESTAMP {
				chunk.hasExtendTimestamp = true
			} else {
				chunk.hasExtendTimestamp = false
				chunk.msgHeader.timestamp += uint64(chunk.msgHeader.timestampDelta)
			}
		case RTMP_FMT_TYPE3:
			// update the timestamp even fmt=3 for first chunk packet. the same with previous.
			if 0 == chunk.msgPayloadSize && !chunk.hasExtendTimestamp {
				chunk.msgHeader.timestamp += uint64(chunk.msgHeader.timestampDelta)
			}
		}

		if err != nil {
			break
		}

		//read extend timestamp
		if chunk.hasExtendTimestamp {
			var buf [4]uint8
			err = rtmp.ExpectBytes(4, buf[:])
			if err != nil {
				break
			}

			extendTimeStamp := binary.BigEndian.Uint32(buf[0:4])

			// always use 31bits timestamp, for some server may use 32bits extended timestamp.
			extendTimeStamp &= 0x7fffffff

			chunkTimeStamp := chunk.msgHeader.timestamp
			if 0 == chunk.msgPayloadSize || 0 == chunkTimeStamp {
				chunk.msgHeader.timestamp = uint64(extendTimeStamp)
			}

			//because of the flv file format is lower 24bits, and higher 8bit is SI32, so timestamp is
			//31bit.
			chunk.msgHeader.timestamp &= 0x7fffffff

		}

		chunk.msgCount++

		//make cache of msg
		if 0 == len(chunk.msgPayload) {
			chunk.msgPayload = make([]uint8, chunk.msgHeader.msgLength)
		}

		//read chunk data
		remainPayloadSize := chunk.msgHeader.msgLength - chunk.msgPayloadSize
		if remainPayloadSize > RTMP_DEFAULT_CHUNK_SIZE {
			remainPayloadSize = RTMP_DEFAULT_CHUNK_SIZE
		}

	}

	return
}
