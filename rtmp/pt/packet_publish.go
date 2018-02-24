package pt

import (
	"fmt"
	"strings"
)

type PublishPacket struct {
	/**
	 * Name of the command, set to “publish”.
	 */
	CommandName string
	/**
	 * Transaction ID set to 0.
	 */
	TransactionId float64
	/**
	 * Command information object does not exist. Set to null type.
	 */
	CommandObject Amf0Object // null
	/**
	 * Name with which the stream is published.
	 */
	StreamName string

	/**
	* Token value, for authentication. it's optional.
	**/
	TokenStr string

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

	pkt.CommandName, err = Amf0ReadString(data, &offset)
	if err != nil {
		return
	}

	if RTMP_AMF0_COMMAND_PUBLISH != pkt.CommandName {
		err = fmt.Errorf("decode publish packet command name is error.actully=%s", pkt.CommandName)
		return
	}

	pkt.TransactionId, err = Amf0ReadNumber(data, &offset)
	if err != nil {
		return
	}

	err = amf0ReadNull(data, &offset)
	if err != nil {
		return
	}

	var streamNameLocal string
	streamNameLocal, err = Amf0ReadString(data, &offset)
	if err != nil {
		return
	}

	i := strings.Index(streamNameLocal, TokenStr)
	if i < 0 {
		pkt.StreamName = streamNameLocal
	} else {
		pkt.StreamName = streamNameLocal[0:i]
		pkt.TokenStr = streamNameLocal[i+len(TokenStr):]
	}

	if uint32(len(data))-offset > 0 {
		pkt.Type, err = Amf0ReadString(data, &offset)
		if err != nil {
			return
		}
	}

	return
}
func (pkt *PublishPacket) Encode() (data []uint8) {
	data = append(data, amf0WriteString(pkt.CommandName)...)
	data = append(data, amf0WriteNumber(pkt.TransactionId)...)
	data = append(data, amf0WriteNull()...)
	data = append(data, amf0WriteString(pkt.StreamName)...)
	data = append(data, amf0WriteString(pkt.Type)...)

	return
}
func (pkt *PublishPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}
func (pkt *PublishPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverStream
}
