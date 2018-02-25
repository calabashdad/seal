package pt

// CallResPacket response for CallPacket
type CallResPacket struct {
	// CommandName Name of the command.
	CommandName string

	// TransactionID ID of the command, to which the response belongs to
	TransactionID float64

	// CommandObject If there exists any command info this is set, else this is set to null type.
	CommandObject interface{}

	// CommandObjectMarker object type marker
	CommandObjectMarker uint8

	//  Response from the method that was called.
	Response interface{}

	// ResponseMarker response type marker
	ResponseMarker uint8
}

// Decode .
func (pkt *CallResPacket) Decode(data []uint8) (err error) {
	var offset uint32

	if pkt.CommandName, err = Amf0ReadString(data, &offset); err != nil {
		return
	}

	if pkt.TransactionID, err = Amf0ReadNumber(data, &offset); err != nil {
		return
	}

	if pkt.CommandObject, err = amf0ReadAny(data, &pkt.CommandObjectMarker, &offset); err != nil {
		return
	}

	maxOffset := uint32(len(data)) - 1
	if maxOffset-offset > 0 {
		pkt.Response, err = amf0ReadAny(data, &pkt.ResponseMarker, &offset)
		if err != nil {
			return
		}

	}

	return
}

// Encode .
func (pkt *CallResPacket) Encode() (data []uint8) {
	data = append(data, amf0WriteString(pkt.CommandName)...)
	data = append(data, amf0WriteNumber(pkt.TransactionID)...)
	if nil != pkt.CommandObject {
		data = append(data, amf0WriteAny(pkt.CommandObject.(Amf0Object))...)
	}

	if nil != pkt.Response {
		data = append(data, amf0WriteAny(pkt.Response.(Amf0Object))...)
	}

	return
}

// GetMessageType .
func (pkt *CallResPacket) GetMessageType() uint8 {
	return RtmpMsgAmf0CommandMessage
}

// GetPreferCsID .
func (pkt *CallResPacket) GetPreferCsID() uint32 {
	return RtmpCidOverConnection
}
