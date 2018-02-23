package pt

import (
	"fmt"
)

/**
* response for SrsConnectAppPacket.
 */
type ConnectResPacket struct {
	/**
	 * _result or _error; indicates whether the response is result or error.
	 */
	CommandName string

	/**
	 * Transaction ID is 1 for call connect responses
	 */
	TransactionId float64

	/**
	 * Name-value pairs that describe the properties(fmsver etc.) of the connection.
	 */
	Props []Amf0Object

	/**
	 * Name-value pairs that describe the response from|the server. ‘code’,
	 * ‘level’, ‘description’ are names of few among such information.
	 */
	Info []Amf0Object
}

func (pkt *ConnectResPacket) Decode(data []uint8) (err error) {
	var offset uint32

	pkt.CommandName, err = Amf0ReadString(data, &offset)
	if err != nil {
		return
	}

	if pkt.CommandName != RTMP_AMF0_COMMAND_RESULT {
		err = fmt.Errorf("decode connect res packet command name is error. actuall name=%s, should be %s",
			pkt.CommandName, RTMP_AMF0_COMMAND_RESULT)
		return
	}

	pkt.TransactionId, err = Amf0ReadNumber(data, &offset)
	if err != nil {
		return
	}

	if pkt.TransactionId != 1.0 {
		err = fmt.Errorf("decode connect res packet transaction id != 1.0.")
		return
	}

	pkt.Props, err = Amf0ReadObject(data, &offset)
	if err != nil {
		return
	}

	pkt.Info, err = Amf0ReadObject(data, &offset)
	if err != nil {
		return
	}

	return
}

func (pkt *ConnectResPacket) Encode() (data []uint8) {

	data = append(data, Amf0WriteString(pkt.CommandName)...)
	data = append(data, Amf0WriteNumber(pkt.TransactionId)...)
	data = append(data, Amf0WriteObject(pkt.Props)...)
	data = append(data, Amf0WriteObject(pkt.Info)...)

	return
}

func (pkt *ConnectResPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}

func (pkt *ConnectResPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverConnection
}
