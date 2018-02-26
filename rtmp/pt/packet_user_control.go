package pt

import (
	"encoding/binary"
	"fmt"
)

// UserControlPacket User Control Message (4)
// for the EventData is 4bytes.
// Stream Begin(=0)              4-bytes stream ID
// Stream EOF(=1)                4-bytes stream ID
// StreamDry(=2)                 4-bytes stream ID
// SetBufferLength(=3)           8-bytes 4bytes stream ID, 4bytes buffer length.
// StreamIsRecorded(=4)          4-bytes stream ID
// PingRequest(=6)               4-bytes timestamp local server time
// PingResponse(=7)              4-bytes timestamp received ping request.
//
// 3.7. User Control message
// +------------------------------+-------------------------
// | Event Type ( 2- bytes ) | Event Data
// +------------------------------+-------------------------
// Figure 5 Pay load for the ‘User Control Message’.
type UserControlPacket struct {

	// Event type is followed by Event data.
	//  @see: SrcPCUCEventType
	EventType uint16
	EventData uint32

	// ExtraData 4bytes if event_type is SetBufferLength; otherwise 0.
	ExtraData uint32
}

// Decode .
func (pkt *UserControlPacket) Decode(data []uint8) (err error) {
	if len(data) < 6 {
		err = fmt.Errorf("decode usercontrol, data len is less than 6. actually is %d", len(data))
		return
	}

	var offset uint32
	pkt.EventType = binary.BigEndian.Uint16(data[offset : offset+2])
	offset += 2

	pkt.EventData = binary.BigEndian.Uint32(data[offset : offset+4])
	offset += 4

	if SrcPCUCSetBufferLength == pkt.EventType {
		if uint32(len(data))-offset < 4 {
			err = fmt.Errorf("decode user control packet extra data, len is not enough < 4")
			return
		}
		pkt.ExtraData = binary.BigEndian.Uint32(data[offset : offset+4])
	}

	return
}

// Encode .
func (pkt *UserControlPacket) Encode() (data []uint8) {

	if SrcPCUCSetBufferLength == pkt.EventType {
		data = make([]uint8, 10)
	} else {
		data = make([]uint8, 6)
	}

	var offset uint32

	binary.BigEndian.PutUint16(data[offset:offset+2], pkt.EventType)
	offset += 2

	binary.BigEndian.PutUint32(data[offset:offset+4], pkt.EventData)
	offset += 4

	if SrcPCUCSetBufferLength == pkt.EventType {
		binary.BigEndian.PutUint32(data[offset:offset+4], pkt.ExtraData)
		offset += 4
	}

	return
}

// GetMessageType .
func (pkt *UserControlPacket) GetMessageType() uint8 {
	return RtmpMsgUserControlMessage
}

// GetPreferCsID .
func (pkt *UserControlPacket) GetPreferCsID() uint32 {
	return RtmpCidProtocolControl
}
