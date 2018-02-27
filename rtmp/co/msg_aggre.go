package co

import (
	"encoding/binary"
	"fmt"
	"log"
	"seal/rtmp/pt"

	"github.com/calabashdad/utiltools"
)

func (rc *RtmpConn) msgAggregate(msg *pt.Message) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	log.Println("aggregate")
	if nil == msg {
		return
	}

	var offset uint32
	maxOffset := len(msg.Payload.Payload) - 1
	if maxOffset < 0 {
		return
	}
	for {
		if offset > uint32(maxOffset) {
			break
		}

		dataType := msg.Payload.Payload[offset]
		offset++

		_ = dataType

		if maxOffset-int(offset) < 3 {
			err = fmt.Errorf("aggregate msg size invalid")
			break
		}
		// data size, 3 bytes, big endian
		dataSize := uint32(msg.Payload.Payload[offset])<<16 + uint32(msg.Payload.Payload[offset+1])<<8 + uint32(msg.Payload.Payload[offset+2])
		if dataSize <= 0 {
			err = fmt.Errorf("aggregate msg size invalid, 0")
			break
		}
		offset += 3

		if maxOffset-int(offset) < 3 {
			err = fmt.Errorf("aggregate msg timestamp invalid")
			break
		}
		timeStamp := uint32(msg.Payload.Payload[offset])<<16 + uint32(msg.Payload.Payload[offset+1])<<8 + uint32(msg.Payload.Payload[offset+2])
		offset += 3

		if maxOffset-int(offset) < 1 {
			err = fmt.Errorf("aggregate msg timeH invalid")
			break
		}
		timeH := uint32(msg.Payload.Payload[offset])
		offset++

		timeStamp |= timeH << 24
		timeStamp &= 0x7FFFFFFF

		if maxOffset-int(offset) < 3 {
			err = fmt.Errorf("aggregate msg stream id invalid")
			break
		}
		streamID := uint32(msg.Payload.Payload[offset])<<16 + uint32(msg.Payload.Payload[offset+1])<<8 + uint32(msg.Payload.Payload[offset+2])
		offset += 3

		if maxOffset-int(offset) < int(dataSize) {
			err = fmt.Errorf("aggregate msg data size, not enough")
			break
		}

		var o pt.Message
		o.Header.MessageType = dataType
		o.Header.PayloadLength = dataSize
		o.Header.TimestampDelta = timeStamp
		o.Header.Timestamp = uint64(timeStamp)
		o.Header.StreamID = streamID
		o.Header.PerferCsid = msg.Header.PerferCsid

		o.Payload.Payload = make([]uint8, dataSize)
		copy(o.Payload.Payload, msg.Payload.Payload[offset:offset+dataSize])
		offset += dataSize

		if maxOffset-int(offset) < 4 {
			err = fmt.Errorf("aggregate msg previous tag size")
			break
		}
		_ = binary.BigEndian.Uint32(msg.Payload.Payload[offset : offset+4])
		offset += 4

		// process has parsed message
		if o.Header.IsAudio() {
			log.Println("aggregate audio msg")
			if err = rc.msgAudio(&o); err != nil {
				break
			}
		} else if o.Header.IsVideo() {
			log.Println("aggregate video msg")
			if err = rc.msgVideo(&o); err != nil {
				break
			}
		}
	}

	return
}
