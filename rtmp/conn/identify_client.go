package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"reflect"
	"seal/rtmp/protocol"
)

func (rc *RtmpConn) IdentifyClient() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	for {
		var msg *protocol.Message

		err = rc.RecvMsg(&msg)
		if err != nil {
			break
		}

		h := msg.Header

		if h.IsAckledgement() || h.IsSetChunkSize() || h.IsWindowAckledgementSize() || h.IsUserControlMessage() {
			continue
		}

		if !h.IsAmf0Command() && !h.IsAmf3Command() {
			continue
		}

		var pkt protocol.Packet

		err = rc.DecodeMsg(&msg, &pkt)
		if err != nil {
			break
		}

		var pktCreateStream *protocol.CreateStreamPacket
		var pktFMLEStart *protocol.FmleStartPacket
		var pktPlay *protocol.PlayPacket
		var pktCallRes *protocol.CallResPacket
		switch reflect.TypeOf(pkt) {
		case reflect.TypeOf(pktCreateStream):
			pktCreateStream = pkt.(*protocol.CreateStreamPacket)
			return rc.identifyCreateStreamClient(pktCreateStream)
		case reflect.TypeOf(pktFMLEStart):
			pktFMLEStart = pkt.(*protocol.FmleStartPacket)
			return
		case reflect.TypeOf(pktPlay):
			pktPlay = pkt.(*protocol.PlayPacket)
			return
		case reflect.TypeOf(pktCallRes):
			pktCallRes = pkt.(*protocol.CallResPacket)
		}
	}

	if err != nil {
		return
	}

	return
}

func (rc *RtmpConn) identifyCreateStreamClient(req *protocol.CreateStreamPacket) (err error) {

	var pkt protocol.CreateStreamResPacket

	pkt.Command_name = protocol.RTMP_AMF0_COMMAND_RESULT
	pkt.Transaction_id = req.Transaction_id
	pkt.Stream_id = 1 //default for the response of create stream.

	err = rc.SendPacket(&pkt, 0)
	if err != nil {
		return
	}

	return
}

func (rc *RtmpConn) identifyFmlePublishClient() (err error) {
	return
}

func (rc *RtmpConn) identifyPlayClient() (err error) {
	return
}
