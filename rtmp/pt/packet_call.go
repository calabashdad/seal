package pt

// CallPacket  the call method of the NetConnection object runs remote procedure
// calls (RPC) at the receiving end. The called RPC name is passed as a parameter to the
// call command
type CallPacket struct {

	// CommandName Name of the remote procedure that is called.
	CommandName string

	// TransactionID If a response is expected we give a transaction Id. Else we pass a value of 0
	TransactionID float64

	// CommandObject If there exists any command info this
	// is set, else this is set to null type.
	CommandObject interface{}

	// CmdObjectType object type marker
	CmdObjectType uint8

	// Arguments Any optional arguments to be provided
	Arguments interface{}

	// ArgumentsType type of Arguments
	ArgumentsType uint8
}

// Decode .
func (pkt *CallPacket) Decode(data []uint8) (err error) {
	var offset uint32

	if pkt.CommandName, err = Amf0ReadString(data, &offset); err != nil {
		return
	}

	if pkt.TransactionID, err = Amf0ReadNumber(data, &offset); err != nil {
		return
	}

	if pkt.CommandObject, err = amf0ReadAny(data, &pkt.CmdObjectType, &offset); err != nil {
		return
	}

	if uint32(len(data))-offset > 0 {
		pkt.Arguments, err = amf0ReadAny(data, &pkt.ArgumentsType, &offset)
		if err != nil {
			return
		}
	}

	return
}

// Encode .
func (pkt *CallPacket) Encode() (data []uint8) {
	data = append(data, amf0WriteString(pkt.CommandName)...)
	data = append(data, amf0WriteNumber(pkt.TransactionID)...)

	if nil != pkt.CommandObject {
		data = append(data, amf0WriteAny(pkt.CommandObject.(Amf0Object))...)
	}

	if nil != pkt.Arguments {
		data = append(data, amf0WriteAny(pkt.Arguments.(Amf0Object))...)
	}

	return
}

// GetMessageType .
func (pkt *CallPacket) GetMessageType() uint8 {
	return RtmpMsgAmf0CommandMessage
}

// GetPreferCsID .
func (pkt *CallPacket) GetPreferCsID() uint32 {
	return RtmpCidOverConnection
}
