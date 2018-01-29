package protocol

/**
* 4.1.1. connect
* The client sends the connect command to the server to request
* connection to a server application instance.
 */
type ConnectPacket struct {
}

func (pkt *ConnectPacket) Decode([]uint8) (err error) {
	return
}
func (pkt *ConnectPacket) Encode() (b []uint8) {
	return
}
func (pkt *ConnectPacket) GetMessageType() uint8 {
	return RTMP_MSG_AMF0CommandMessage
}
func (pkt *ConnectPacket) GetPreferCsId() uint32 {
	return RTMP_CID_OverConnection
}
