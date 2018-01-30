package pt

import (
	"fmt"
)

type PublishPacket struct {
	/**
	 * Name of the command, set to “publish”.
	 */
	Command_name string
	/**
	 * Transaction ID set to 0.
	 */
	Transaction_id float64
	/**
	 * Command information object does not exist. Set to null type.
	 */
	Command_object Amf0Object // null
	/**
	 * Name with which the stream is published.
	 */
	Stream_name string
	/**
	 * Type of publishing. Set to “live”, “record”, or “append”.
	 *   record: The stream is published and the data is recorded to a new file.The file
	 *           is stored on the server in a subdirectory within the directory that
	 *           contains the server application. If the file already exists, it is
	 *           overwritten.
	 *   append: The stream is published and the data is appended to a file. If no file
	 *           is found, it is created.
	 *   live: Live data is published without recording it in a file.
	 * @remark, SRS only support live.
	 * @remark, optional, default to live.
	 */
	Type string
}

func (pkt *PublishPacket) Decode(data []uint8) (err error) {
	var offset uint32

	err, pkt.Command_name = Amf0ReadString(data, &offset)
	if err != nil {
		return
	}

	if RTMP_AMF0_COMMAND_PUBLISH != pkt.Command_name {
		err = fmt.Errorf("decode publish packet command name is error.actully=%s", pkt.Command_name)
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

	err, pkt.Stream_name = Amf0ReadString(data, &offset)
	if err != nil {
		return
	}

	if uint32(len(data))-offset > 0 {
		err, pkt.Type = Amf0ReadString(data, &offset)
		if err != nil {
			return
		}
	}

	return
}
func (pkt *PublishPacket) Encode() (data []uint8) {
	data = append(data, Amf0WriteString(pkt.Command_name)...)
	data = append(data, Amf0WriteNumber(pkt.Transaction_id)...)
	data = append(data, Amf0WriteNull()...)
	data = append(data, Amf0WriteString(pkt.Stream_name)...)
	data = append(data, Amf0WriteString(pkt.Type)...)
	
	return
}
func (pkt *PublishPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}
func (pkt *PublishPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverStream
}
