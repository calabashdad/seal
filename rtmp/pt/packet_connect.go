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
	CommandName string
	/**
	 * Always set to 1.
	 */
	TransactionId float64
	/**
	 * Command information object which has the name-value pairs.
	 * @remark: alloc in packet constructor, user can directly use it,
	 *       user should never alloc it again which will cause memory leak.
	 */
	CommandObject []Amf0Object
	/**
	 * Any optional information
	 */
	Args []Amf0Object
}

func (pkt *ConnectPacket) GetObjectProperty(name string) (value interface{}) {

	for _, v := range pkt.CommandObject {
		if name == v.PropertyName {
			return v.Value
		}
	}

	return
}

func (pkt *ConnectPacket) Decode(data []uint8) (err error) {
	var offset uint32

	err, pkt.CommandName = Amf0ReadString(data, &offset)
	if err != nil {
		return
	}

	if RTMP_AMF0_COMMAND_CONNECT != pkt.CommandName {
		err = fmt.Errorf("decode connect packet, command name is not connect.")
		return
	}

	err, pkt.TransactionId = Amf0ReadNumber(data, &offset)
	if err != nil {
		return
	}

	if 1.0 != pkt.TransactionId {
		err = fmt.Errorf("decode conenct packet, transaction id is not 1.0")
		return
	}

	err, pkt.CommandObject = Amf0ReadObject(data, &offset)
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
	data = append(data, Amf0WriteString(pkt.CommandName)...)
	data = append(data, Amf0WriteNumber(pkt.TransactionId)...)
	data = append(data, Amf0WriteObject(pkt.CommandObject)...)
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
