package pt

import "fmt"

// ConnectPacket The client sends the connect command to the server to request
// connection to a server application instance.
type ConnectPacket struct {

	// CommandName Name of the command. Set to “connect”.
	CommandName string

	// TransactionID Always set to 1.
	TransactionID float64

	// CommandObject Command information object which has the name-value pairs.
	CommandObject []Amf0Object

	// Args Any optional information
	Args []Amf0Object
}

// GetObjectProperty get object property in connect packet
func (pkt *ConnectPacket) GetObjectProperty(name string) (value interface{}) {

	for _, v := range pkt.CommandObject {
		if name == v.propertyName {
			return v.value
		}
	}

	return
}

// Decode .
func (pkt *ConnectPacket) Decode(data []uint8) (err error) {
	var offset uint32

	if pkt.CommandName, err = Amf0ReadString(data, &offset); err != nil {
		return
	}

	if RtmpAmf0CommandConnect != pkt.CommandName {
		err = fmt.Errorf("decode connect packet, command name is not connect")
		return
	}

	if pkt.TransactionID, err = Amf0ReadNumber(data, &offset); err != nil {
		return
	}

	if 1.0 != pkt.TransactionID {
		err = fmt.Errorf("decode conenct packet, transaction id is not 1.0")
		return
	}

	if pkt.CommandObject, err = amf0ReadObject(data, &offset); err != nil {
		return
	}

	if uint32(len(data))-offset > 0 {
		var marker uint8
		var v interface{}
		v, err = amf0ReadAny(data, &marker, &offset)
		if err != nil {
			return
		}

		if RtmpAmf0Object == marker {
			pkt.Args = v.([]Amf0Object)
		}
	}

	return
}

// Encode .
func (pkt *ConnectPacket) Encode() (data []uint8) {
	data = append(data, amf0WriteString(pkt.CommandName)...)
	data = append(data, amf0WriteNumber(pkt.TransactionID)...)
	data = append(data, amf0WriteObject(pkt.CommandObject)...)
	if len(pkt.Args) > 0 {
		data = append(data, amf0WriteObject(pkt.Args)...)
	}

	return
}

// GetMessageType .
func (pkt *ConnectPacket) GetMessageType() uint8 {
	return RtmpMsgAmf0CommandMessage
}

// GetPreferCsID .
func (pkt *ConnectPacket) GetPreferCsID() uint32 {
	return RtmpCidOverConnection
}
