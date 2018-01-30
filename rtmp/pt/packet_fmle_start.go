package pt

import (
	"fmt"
	"strings"
)

/**
* FMLE start publish: ReleaseStream/PublishStream
 */
type FmleStartPacket struct {
	/**
	 * Name of the command
	 */
	Command_name string
	/**
	 * the transaction ID to get the response.
	 */
	Transaction_id float64
	/**
	 * If there exists any command info this is set, else this is set to null type.
	 */
	Command_object Amf0Object // null
	/**
	 * the stream name to start publish or release.
	 */
	StreamName string

	/**
	* Token value, for authentication. it's optional.
	**/
	TokenStr string
}

func (pkt *FmleStartPacket) Decode(data []uint8) (err error) {
	var offset uint32

	err, pkt.Command_name = Amf0ReadString(data, &offset)
	if err != nil {
		return
	}

	if RTMP_AMF0_COMMAND_RELEASE_STREAM != pkt.Command_name &&
		RTMP_AMF0_COMMAND_FC_PUBLISH != pkt.Command_name &&
		RTMP_AMF0_COMMAND_UNPUBLISH != pkt.Command_name {
		err = fmt.Errorf("decode fmle start packet error, command name is error.actully=", pkt.Command_name)
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

	var streamNameLocal string
	err, streamNameLocal = Amf0ReadString(data, &offset)
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

	return
}
func (pkt *FmleStartPacket) Encode() (data []uint8) {
	data = append(data, Amf0WriteString(pkt.Command_name)...)
	data = append(data, Amf0WriteNumber(pkt.Transaction_id)...)
	data = append(data, Amf0WriteNull()...)
	data = append(data, Amf0WriteString(pkt.StreamName)...)

	return
}
func (pkt *FmleStartPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}
func (pkt *FmleStartPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverConnection
}
