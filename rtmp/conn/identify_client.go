package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"reflect"
	"seal/rtmp/pt"
)

func (rc *RtmpConn) IdentifyClient() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	for {
		var msg *pt.Message

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

		var pkt pt.Packet

		err = rc.DecodeMsg(&msg, &pkt)
		if err != nil {
			break
		}

		var pktCreateStream *pt.CreateStreamPacket
		var pktFMLEStart *pt.FmleStartPacket
		var pktPlay *pt.PlayPacket
		var pktCallRes *pt.CallResPacket
		switch reflect.TypeOf(pkt) {
		case reflect.TypeOf(pktCreateStream):
			pktCreateStream = pkt.(*pt.CreateStreamPacket)
			return rc.identifyCreateStreamClient(pktCreateStream)
		case reflect.TypeOf(pktFMLEStart):
			pktFMLEStart = pkt.(*pt.FmleStartPacket)
			return rc.identifyFmlePublishClient(pktFMLEStart)
		case reflect.TypeOf(pktPlay):
			pktPlay = pkt.(*pt.PlayPacket)
			return rc.identifyPlayClient(pktPlay)
		case reflect.TypeOf(pktCallRes):
			pktCallRes = pkt.(*pt.CallResPacket)
		}
	}

	if err != nil {
		return
	}

	return
}

func (rc *RtmpConn) identifyCreateStreamClient(req *pt.CreateStreamPacket) (err error) {

	var pkt pt.CreateStreamResPacket

	pkt.Command_name = pt.RTMP_AMF0_COMMAND_RESULT
	pkt.Transaction_id = req.Transaction_id
	pkt.Stream_id = 1 //default for the response of create stream.

	err = rc.SendPacket(&pkt, 0)
	if err != nil {
		return
	}

	for {
		var msg *pt.Message

		err = rc.RecvMsg(&msg)
		if err != nil {
			break
		}

		if msg.Header.IsAckledgement() || msg.Header.IsSetChunkSize() || msg.Header.IsWindowAckledgementSize() || msg.Header.IsUserControlMessage() {
			continue
		}

		if !msg.Header.IsAmf0Command() && !msg.Header.IsAmf3Command() {
			continue
		}

		var pktLocal pt.Packet
		err = rc.DecodeMsg(&msg, &pktLocal)
		if err != nil {
			break
		}

		var pktPlayPacket *pt.PlayPacket
		var pktPublishPacket *pt.PublishPacket
		var pktCreateStreamPacket *pt.CreateStreamPacket
		switch reflect.TypeOf(pktLocal) {
		case reflect.TypeOf(pktPlayPacket):
			pktPlayPacket = pktLocal.(*pt.PlayPacket)
			return rc.identifyPlayClient(pktPlayPacket)
		case reflect.TypeOf(pktPublishPacket):
			pktPublishPacket = pktLocal.(*pt.PublishPacket)
			return rc.identifyFlashPublishClient(pktPublishPacket)
		case reflect.TypeOf(pktCreateStreamPacket):
			pktCreateStreamPacket = pktLocal.(*pt.CreateStreamPacket)
			return rc.IdentifyClient()
		}
	}

	if err != nil {
		return
	}

	return
}

//publish role.
func (rc *RtmpConn) identifyFmlePublishClient(pktPublishPacket *pt.FmleStartPacket) (err error) {

	if nil == pktPublishPacket {
		return
	}

	rc.Role = RtmpRoleFMLEPublisher
	rc.StreamName = pktPublishPacket.StreamName

	var pktRes pt.FmleStartResPacket
	pktRes.Transaction_id = pktPublishPacket.Transaction_id

	err = rc.SendPacket(&pktRes, 0)
	if err != nil {
		return
	}

	return
}

func (rc *RtmpConn) identifyFlashPublishClient(pktPublishPacket *pt.PublishPacket) (err error) {

	if nil == pktPublishPacket {
		return
	}

	rc.Role = RtmpRoleFMLEPublisher
	rc.StreamName = pktPublishPacket.StreamName

	return
}

//play role.
func (rc *RtmpConn) identifyPlayClient(pktPlayPacket *pt.PlayPacket) (err error) {

	if nil == pktPlayPacket {
		return
	}

	rc.Role = RtmpRolePlayer
	rc.StreamName = pktPlayPacket.StreamName
	rc.Duration = pktPlayPacket.Duration

	return
}
