package protocol

import (
	"fmt"
)

/**
* 4.1.3. createStream
* The client sends this command to the server to create a logical
* channel for message communication The publishing of audio, video, and
* metadata is carried out over stream channel created using the
* createStream command.
 */
type CreateStreamPacket struct {
	/**
	 * Name of the command. Set to “createStream”.
	 */
	Command_name string
	/**
	 * Transaction ID of the command.
	 */
	Transaction_id float64
	/**
	 * If there exists any command info this is set, else this is set to null type.
	 */
	Command_object Amf0Object // null
}

func (pkt *CreateStreamPacket) Decode(data []uint8) (err error) {
	var offset uint32

	err, pkt.Command_name = Amf0ReadString(data, &offset)
	if err != nil {
		return
	}

	if RTMP_AMF0_COMMAND_CREATE_STREAM != pkt.Command_name {
		err = fmt.Errorf("decode create stream packet, command name is wrong.", pkt.Command_name)
		return
	}

	err, pkt.Transaction_id = Amf0ReadNumber(data, &offset)
	if err != nil {
		return
	}

	err = Amf0ReadNull(data, &offset)
	if err != nil {
		return
	}

	return
}
func (pkt *CreateStreamPacket) Encode() (data []uint8) {
	data = append(data, Amf0WriteString(pkt.Command_name)...)
	data = append(data, Amf0WriteNumber(pkt.Transaction_id)...)
	data = append(data, Amf0WriteNull()...)
	
	return
}
func (pkt *CreateStreamPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}
func (pkt *CreateStreamPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverConnection
}
