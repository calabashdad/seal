package pt

import (
	"encoding/binary"
	"fmt"
)

/**
* 5.4. User Control Message (4)
*
* for the EventData is 4bytes.
* Stream Begin(=0)              4-bytes stream ID
* Stream EOF(=1)                4-bytes stream ID
* StreamDry(=2)                 4-bytes stream ID
* SetBufferLength(=3)           8-bytes 4bytes stream ID, 4bytes buffer length.
* StreamIsRecorded(=4)          4-bytes stream ID
* PingRequest(=6)               4-bytes timestamp local server time
* PingResponse(=7)              4-bytes timestamp received ping request.
*
* 3.7. User Control message
* +------------------------------+-------------------------
* | Event Type ( 2- bytes ) | Event Data
* +------------------------------+-------------------------
* Figure 5 Pay load for the ‘User Control Message’.
 */
type UserControlPacket struct {
	/**
	 * Event type is followed by Event data.
	 * @see: SrcPCUCEventType
	 */
	Event_type uint16
	Event_data uint32
	/**
	 * 4bytes if event_type is SetBufferLength; otherwise 0.
	 */
	Extra_data uint32
}

func (pkt *UserControlPacket) Decode(data []uint8) (err error) {
	if len(data) < 6 {
		err = fmt.Errorf("decode usercontrol, data len is less than 6. actually is ", len(data))
		return
	}

	var offset uint32
	pkt.Event_type = binary.BigEndian.Uint16(data[offset : offset+2])
	offset += 2

	pkt.Event_data = binary.BigEndian.Uint32(data[offset : offset+4])
	offset += 4

	if SrcPCUCSetBufferLength == pkt.Event_type {
		if uint32(len(data))-offset < 4 {
			err = fmt.Errorf("decode user control packet extra data, len is not enough. < 4.")
			return
		}
		pkt.Extra_data = binary.BigEndian.Uint32(data[offset : offset+4])
	}

	return
}
func (pkt *UserControlPacket) Encode() (data []uint8) {

	if SrcPCUCSetBufferLength == pkt.Event_type {
		data = make([]uint8, 10)
	} else {
		data = make([]uint8, 6)
	}

	var offset uint32

	binary.BigEndian.PutUint16(data[offset:offset+2], pkt.Event_type)
	offset += 2

	binary.BigEndian.PutUint32(data[offset:offset+4], pkt.Event_data)
	offset += 4

	if SrcPCUCSetBufferLength == pkt.Event_type {
		binary.BigEndian.PutUint32(data[offset:offset+4], pkt.Extra_data)
		offset += 4
	}

	return
}
func (pkt *UserControlPacket) GetMessageType() uint8 {
	return RTMP_MSG_UserControlMessage
}
func (pkt *UserControlPacket) GetPreferCsId() uint32 {
	return RTMP_CID_ProtocolControl
}
