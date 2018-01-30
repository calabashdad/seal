package pt

import "fmt"

/**
* 4.1.1. connect
* The client sends the connect command to the server to request
* connection to a server application instance.
 */
type ConnectPacket struct {
	/**
	 * Name of the command. Set to “connect”.
	 */
	Command_name string
	/**
	 * Always set to 1.
	 */
	Transaction_id float64
	/**
	 * Command information object which has the name-value pairs.
	 * @remark: alloc in packet constructor, user can directly use it,
	 *       user should never alloc it again which will cause memory leak.
	 */
	Command_object []Amf0Object
	/**
	 * Any optional information
	 */
	Args []Amf0Object
}

func (pkt *ConnectPacket) Decode(data []uint8) (err error) {
	var offset uint32

	err, pkt.Command_name = Amf0ReadString(data, &offset)
	if err != nil {
		return
	}

	if RTMP_AMF0_COMMAND_CONNECT != pkt.Command_name {
		err = fmt.Errorf("decode connect packet, command name is not connect.")
		return
	}

	err, pkt.Transaction_id = Amf0ReadNumber(data, &offset)
	if err != nil {
		return
	}

	if 1.0 != pkt.Transaction_id {
		err = fmt.Errorf("decode conenct packet, transaction id is not 1.0")
		return
	}

	err, pkt.Command_object = Amf0ReadObject(data, &offset)
	if err != nil {
		return
	}

	if uint32(len(data))-offset > 0 {
		var marker uint8
		var v interface{}
		err, v = Amf0ReadAny(data, &marker, &offset)
		if err != nil {
			return
		}

		if RTMP_AMF0_Object == marker {
			pkt.Args = v.([]Amf0Object)
		}
	}

	return
}
func (pkt *ConnectPacket) Encode() (data []uint8) {
	data = append(data, Amf0WriteString(pkt.Command_name)...)
	data = append(data, Amf0WriteNumber(pkt.Transaction_id)...)
	data = append(data, Amf0WriteObject(pkt.Command_object)...)
	if len(pkt.Args) > 0 {
		data = append(data, Amf0WriteObject(pkt.Args)...)
	}

	return
}
func (pkt *ConnectPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}
func (pkt *ConnectPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverConnection
}
