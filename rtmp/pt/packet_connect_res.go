package pt

import (
	"fmt"
)

// ConnectResPacket response for SrsConnectAppPacket.
type ConnectResPacket struct {

	// CommandName _result or _error; indicates whether the response is result or error.
	CommandName string

	// Transaction ID is 1 for call connect responses
	TransactionID float64

	// Props Name-value pairs that describe the properties(fmsver etc.) of the connection.
	Props []Amf0Object

	// Info Name-value pairs that describe the response from|the server. ‘code’,
	// ‘level’, ‘description’ are names of few among such information.
	Info []Amf0Object
}

// Decode .
func (pkt *ConnectResPacket) Decode(data []uint8) (err error) {
	var offset uint32

	if pkt.CommandName, err = Amf0ReadString(data, &offset); err != nil {
		return
	}

	if pkt.CommandName != RTMP_AMF0_COMMAND_RESULT {
		err = fmt.Errorf("decode connect res packet command name is error. actuall name=%s, should be %s",
			pkt.CommandName, RTMP_AMF0_COMMAND_RESULT)
		return
	}

	if pkt.TransactionID, err = Amf0ReadNumber(data, &offset); err != nil {
		return
	}

	if pkt.TransactionID != 1.0 {
		err = fmt.Errorf("decode connect res packet transaction id != 1.0")
		return
	}

	if pkt.Props, err = amf0ReadObject(data, &offset); err != nil {
		return
	}

	if pkt.Info, err = amf0ReadObject(data, &offset); err != nil {
		return
	}

	return
}

// Encode .
func (pkt *ConnectResPacket) Encode() (data []uint8) {

	data = append(data, amf0WriteString(pkt.CommandName)...)
	data = append(data, amf0WriteNumber(pkt.TransactionID)...)
	data = append(data, amf0WriteObject(pkt.Props)...)
	data = append(data, amf0WriteObject(pkt.Info)...)

	return
}

// GetMessageType .
func (pkt *ConnectResPacket) GetMessageType() uint8 {
	return RtmpMsgAmf0CommandMessage
}

// GetPreferCsID .
func (pkt *ConnectResPacket) GetPreferCsID() uint32 {
	return RtmpCidOverConnection
}

// AddProsObj add object to pros
func (pkt *ConnectResPacket) AddProsObj(obj *Amf0Object) {
	pkt.Props = append(pkt.Props, *obj)
}
