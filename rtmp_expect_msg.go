package main

import (
	"UtilsTools/identify_panic"
	"encoding/binary"
	"fmt"
	"log"
)

type MessageStream struct {
	header struct {
		timestampDelta uint32
		timestamp      uint64
		length         uint32
		typeId         uint8
		streamId       uint32
		preferCsId     uint32
	}

	payload []uint8
}

type ChunkStream struct {
	chunkFmt uint8
	csId     uint32
	//msg count of this chunk
	msgCount      uint64
	msgHeaderSize uint32

	msg MessageStream

	hasExtendTimestamp bool
	extendTimeStamp    uint32

	//decode message this time. when finished, will be reset to 0.
	payloadSizeTmp uint32

	decodeResultType string
	decodeResult     interface{}
}

func (rtmp *RtmpSession) RecvMsg() (err error, chunk *ChunkStream) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}
	}()

	//expect msg.
	for {
		//read basic header
		var buf [3]uint8

		err = rtmp.ExpectBytes(1, buf[:1])
		if err != nil {
			return
		}

		chunk_fmt := (buf[0] & 0xc0) >> 6
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

		_, ok := rtmp.chunks[csid]
		if !ok {
			chunk = &ChunkStream{}

			rtmp.chunks[csid] = chunk
		}

		chunk.chunkFmt = chunk_fmt
		chunk.csId = csid
		chunk.msg.header.preferCsId = csid

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

		if chunk.payloadSizeTmp > 0 && RTMP_FMT_TYPE0 == chunk.chunkFmt {
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
			chunk.msg.header.timestampDelta = uint32(msgHeader[0])<<16 + uint32(msgHeader[1])<<8 + uint32(msgHeader[2])
			if chunk.msg.header.timestampDelta >= RTMP_EXTENDED_TIMESTAMP {
				chunk.hasExtendTimestamp = true
			} else {
				chunk.hasExtendTimestamp = false
				// For a type-0 chunk, the absolute timestamp of the message is sent here.
				chunk.msg.header.timestamp = uint64(chunk.msg.header.timestampDelta)
			}

			payloadLength := uint32(msgHeader[3])<<16 + uint32(msgHeader[4])<<8 + uint32(msgHeader[5])
			if chunk.payloadSizeTmp > 0 && payloadLength != chunk.msg.header.length {
				err = fmt.Errorf("RTMP_FMT_TYPE0: msg has in chunk, msg size can not change.")
			}

			chunk.msg.header.length = payloadLength

			chunk.msg.header.typeId = msgHeader[6]

			chunk.msg.header.streamId = binary.LittleEndian.Uint32(msgHeader[7:11])

		case RTMP_FMT_TYPE1:
			chunk.msg.header.timestampDelta = uint32(msgHeader[0])<<16 + uint32(msgHeader[1])<<8 + uint32(msgHeader[2])
			if chunk.msg.header.timestampDelta >= RTMP_EXTENDED_TIMESTAMP {
				chunk.hasExtendTimestamp = true
			} else {
				chunk.hasExtendTimestamp = false
				chunk.msg.header.timestamp += uint64(chunk.msg.header.timestampDelta)
			}

			payloadLength := uint32(msgHeader[3])<<16 + uint32(msgHeader[4])<<8 + uint32(msgHeader[5])
			if chunk.payloadSizeTmp > 0 && payloadLength != chunk.msg.header.length {
				err = fmt.Errorf("RTMP_FMT_TYPE1: msg has in chunk, msg size can not change.")
			}

			chunk.msg.header.length = payloadLength

			chunk.msg.header.typeId = msgHeader[6]

		case RTMP_FMT_TYPE2:
			chunk.msg.header.timestampDelta = uint32(msgHeader[0])<<16 + uint32(msgHeader[1])<<8 + uint32(msgHeader[2])
			if chunk.msg.header.timestampDelta >= RTMP_EXTENDED_TIMESTAMP {
				chunk.hasExtendTimestamp = true
			} else {
				chunk.hasExtendTimestamp = false
				chunk.msg.header.timestamp += uint64(chunk.msg.header.timestampDelta)
			}
		case RTMP_FMT_TYPE3:
			// update the timestamp even fmt=3 for first chunk packet. the same with previous.
			if 0 == chunk.payloadSizeTmp && !chunk.hasExtendTimestamp {
				chunk.msg.header.timestamp += uint64(chunk.msg.header.timestampDelta)
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

			chunkTimeStamp := chunk.msg.header.timestamp
			if 0 == chunk.payloadSizeTmp || 0 == chunkTimeStamp {
				chunk.msg.header.timestamp = uint64(extendTimeStamp)
			}

			//because of the flv file format is lower 24bits, and higher 8bit is SI32, so timestamp is
			//31bit.
			chunk.msg.header.timestamp &= 0x7fffffff

		}

		chunk.msgCount++

		//make cache of msg
		if uint32(len(chunk.msg.payload)) < chunk.msg.header.length {
			chunk.msg.payload = make([]uint8, chunk.msg.header.length)
		}

		//read chunk data
		remainPayloadSize := chunk.msg.header.length - chunk.payloadSizeTmp

		if remainPayloadSize >= rtmp.chunkSize {
			remainPayloadSize = rtmp.chunkSize
		}

		err = rtmp.ExpectBytes(remainPayloadSize, chunk.msg.payload[chunk.payloadSizeTmp:chunk.payloadSizeTmp+remainPayloadSize])
		if err != nil {
			break
		} else {
			chunk.payloadSizeTmp += remainPayloadSize
			if chunk.payloadSizeTmp == chunk.msg.header.length {

				//has recv entire rtmp message.
				//reset the payload size this time, the message actually size is header length, this chunk can reuse by a new csid.
				chunk.payloadSizeTmp = 0

				break
			}
		}

	}

	if err != nil {
		return
	}

	err = rtmp.EstimateNeedSendAcknowlegement(chunk)
	if err != nil {
		return
	}
	return
}

func (rtmp *RtmpSession) EstimateNeedSendAcknowlegement(chunk *ChunkStream) (err error) {
	if (rtmp.ackWindow.ackWindowSize > 0) && (rtmp.recvBytesSum-rtmp.ackWindow.hasAckedSize > uint64(rtmp.ackWindow.ackWindowSize)) {

		err = rtmp.CommonMsgResponseWindowAcknowledgement(chunk, uint32(rtmp.recvBytesSum))
		if err != nil {
			return
		}

		rtmp.ackWindow.hasAckedSize = rtmp.recvBytesSum
	}

	return
}
