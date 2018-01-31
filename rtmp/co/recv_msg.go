package co

import (
	"UtilsTools/identify_panic"
	"encoding/binary"
	"fmt"
	"log"
	"seal/rtmp/pt"
)

//RecvMsg recv whole msg and quit when got an entire msg, not handle it at all.
func (rc *RtmpConn) RecvMsg(chunkStreamID *uint32) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	for {
		//read basic header
		var buf [3]uint8
		err = rc.TcpConn.ExpectBytesFull(buf[:1], 1)
		if err != nil {
			return
		}

		chunkFmt := (buf[0] & 0xc0) >> 6
		csid := uint32(buf[0] & 0x3f)

		switch csid {
		case 0:
			//csId 2 bytes. 64-319
			err = rc.TcpConn.ExpectBytesFull(buf[1:2], 1)
			csid = uint32(64 + buf[1])
		case 1:
			//csId 3 bytes. 64-65599
			err = rc.TcpConn.ExpectBytesFull(buf[1:3], 2)
			csid = uint32(64) + uint32(buf[1]) + uint32(buf[2])*256
		}

		if err != nil {
			break
		}

		chunk, ok := rc.ChunkStreams[csid]
		if !ok {
			chunk = &pt.ChunkStream{}
			rc.ChunkStreams[csid] = chunk
		}

		chunk.Fmt = chunkFmt
		chunk.CsId = csid
		chunk.Msg.Header.PerferCsid = csid

		//read message header
		if 0 == chunk.MsgCount && chunk.Fmt != pt.RTMP_FMT_TYPE0 {
			if pt.RTMP_CID_ProtocolControl == chunk.CsId && pt.RTMP_FMT_TYPE1 == chunk.Fmt {
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

		if chunk.Msg.SizeTmp > 0 && pt.RTMP_FMT_TYPE0 == chunk.Fmt {
			err = fmt.Errorf("when msg count > 0, chunk fmt is not allowed to be RTMP_FMT_TYPE0.")
			break
		}

		var msgHeaderSize uint32

		switch chunk.Fmt {
		case pt.RTMP_FMT_TYPE0:
			msgHeaderSize = 11
		case pt.RTMP_FMT_TYPE1:
			msgHeaderSize = 7
		case pt.RTMP_FMT_TYPE2:
			msgHeaderSize = 3
		case pt.RTMP_FMT_TYPE3:
			msgHeaderSize = 0
		}

		var msgHeader [11]uint8 //max is 11
		err = rc.TcpConn.ExpectBytesFull(msgHeader[:], msgHeaderSize)
		if err != nil {
			break
		}

		//parse msg header
		//*   3bytes: timestamp delta,    fmt=0,1,2
		//*   3bytes: payload length,     fmt=0,1
		//*   1bytes: message type,       fmt=0,1
		//*   4bytes: stream id,          fmt=0
		switch chunk.Fmt {
		case pt.RTMP_FMT_TYPE0:
			chunk.Msg.Header.TimestampDelta = uint32(msgHeader[0])<<16 + uint32(msgHeader[1])<<8 + uint32(msgHeader[2])
			if chunk.Msg.Header.TimestampDelta >= pt.RTMP_EXTENDED_TIMESTAMP {
				chunk.HasExtendedTimestamp = true
			} else {
				chunk.HasExtendedTimestamp = false
				// For a type-0 chunk, the absolute timestamp of the message is sent here.
				chunk.Msg.Header.Timestamp = uint64(chunk.Msg.Header.TimestampDelta)
			}

			payloadLength := uint32(msgHeader[3])<<16 + uint32(msgHeader[4])<<8 + uint32(msgHeader[5])
			if chunk.Msg.SizeTmp > 0 && payloadLength != chunk.Msg.Header.PayloadLength {
				err = fmt.Errorf("RTMP_FMT_TYPE0: msg has in chunk, msg size can not change.")
				break
			}

			chunk.Msg.Header.PayloadLength = payloadLength
			chunk.Msg.Header.MessageType = msgHeader[6]
			chunk.Msg.Header.StreamId = binary.LittleEndian.Uint32(msgHeader[7:11])

		case pt.RTMP_FMT_TYPE1:
			chunk.Msg.Header.TimestampDelta = uint32(msgHeader[0])<<16 + uint32(msgHeader[1])<<8 + uint32(msgHeader[2])
			if chunk.Msg.Header.TimestampDelta >= pt.RTMP_EXTENDED_TIMESTAMP {
				chunk.HasExtendedTimestamp = true
			} else {
				chunk.HasExtendedTimestamp = false
				chunk.Msg.Header.Timestamp += uint64(chunk.Msg.Header.TimestampDelta)
			}

			payloadLength := uint32(msgHeader[3])<<16 + uint32(msgHeader[4])<<8 + uint32(msgHeader[5])
			if chunk.Msg.SizeTmp > 0 && payloadLength != chunk.Msg.Header.PayloadLength {
				err = fmt.Errorf("RTMP_FMT_TYPE1: msg has in chunk, msg size can not change.")
				break
			}

			chunk.Msg.Header.PayloadLength = payloadLength
			chunk.Msg.Header.MessageType = msgHeader[6]

		case pt.RTMP_FMT_TYPE2:
			chunk.Msg.Header.TimestampDelta = uint32(msgHeader[0])<<16 + uint32(msgHeader[1])<<8 + uint32(msgHeader[2])
			if chunk.Msg.Header.TimestampDelta >= pt.RTMP_EXTENDED_TIMESTAMP {
				chunk.HasExtendedTimestamp = true
			} else {
				chunk.HasExtendedTimestamp = false
				chunk.Msg.Header.Timestamp += uint64(chunk.Msg.Header.TimestampDelta)
			}
		case pt.RTMP_FMT_TYPE3:
			// update the timestamp even fmt=3 for first chunk packet. the same with previous.
			if 0 == chunk.Msg.SizeTmp && !chunk.HasExtendedTimestamp {
				chunk.Msg.Header.Timestamp += uint64(chunk.Msg.Header.TimestampDelta)
			}
		}

		if err != nil {
			break
		}

		//read extend timestamp
		if chunk.HasExtendedTimestamp {
			var buf [4]uint8
			err = rc.TcpConn.ExpectBytesFull(buf[:], 4)
			if err != nil {
				break
			}

			extendTimeStamp := binary.BigEndian.Uint32(buf[0:4])

			// always use 31bits timestamp, for some server may use 32bits extended timestamp.
			extendTimeStamp &= 0x7fffffff

			chunkTimeStamp := chunk.Msg.Header.Timestamp
			if 0 == chunk.Msg.SizeTmp || 0 == chunkTimeStamp {
				chunk.Msg.Header.Timestamp = uint64(extendTimeStamp)
			}

			//because of the flv file format is lower 24bits, and higher 8bit is SI32, so timestamp is
			//31bit.
			chunk.Msg.Header.Timestamp &= 0x7fffffff

		}

		chunk.MsgCount++

		//make cache of msg
		if uint32(len(chunk.Msg.Payload)) < chunk.Msg.Header.PayloadLength {
			chunk.Msg.Payload = make([]uint8, chunk.Msg.Header.PayloadLength)
		}

		//read chunk data
		remainPayloadSize := chunk.Msg.Header.PayloadLength - chunk.Msg.SizeTmp

		if remainPayloadSize >= rc.InChunkSize {
			remainPayloadSize = rc.InChunkSize
		}

		err = rc.TcpConn.ExpectBytesFull(chunk.Msg.Payload[chunk.Msg.SizeTmp:chunk.Msg.SizeTmp+remainPayloadSize], remainPayloadSize)
		if err != nil {
			break
		} else {
			chunk.Msg.SizeTmp += remainPayloadSize
			if chunk.Msg.SizeTmp == chunk.Msg.Header.PayloadLength {

				*chunkStreamID = csid
				//has recv entire rtmp message.
				//reset the payload size this time, the message actually size is header length, this chunk can reuse by a new csid.
				chunk.Msg.SizeTmp = 0

				break
			}
		}

	}

	if err != nil {
		return
	}

	return
}
