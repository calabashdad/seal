package conn

import (
	"UtilsTools/identify_panic"
	"fmt"
	"log"
	"seal/rtmp/pt"
)

func (rc *RtmpConn) DecodeMsg(msg **pt.Message, pkt *pt.Packet) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	if nil == msg {
		log.Println("nil when decode msg")
		return
	}

	//offset of parsed msg payload.
	var offset uint32

	if (*msg).Header.IsAmf0Command() || (*msg).Header.IsAmf3Command() ||
		(*msg).Header.IsAmf0Data() || (*msg).Header.IsAmf3Data() {

		if (*msg).Header.IsAmf3Command() && len((*msg).Payload) > 1 {
			offset += 1
		}

		var command string
		err, command = pt.Amf0ReadString((*msg).Payload, &offset)
		if err != nil {
			log.Println("read command failed when decode msg.")
			return
		}

		switch command {
		//command: result or error.
		case pt.RTMP_AMF0_COMMAND_RESULT, pt.RTMP_AMF0_COMMAND_ERROR:
			if true {
				var transaction_id float64
				err, transaction_id = pt.Amf0ReadNumber((*msg).Payload, &offset)
				if err != nil {
					log.Println("read transaction id failed when decode msg.")
					return
				}

				req_command_name := rc.Requests[transaction_id]
				if 0 == len(req_command_name) {
					err = fmt.Errorf("can not find request command name.")
					return
				}

				switch req_command_name {
				case pt.RTMP_AMF0_COMMAND_CONNECT:
					*pkt = &pt.ConnectResPacket{}
					err = (*pkt).Decode((*msg).Payload)

				case pt.RTMP_AMF0_COMMAND_CREATE_STREAM:
					(*pkt) = &pt.CreateStreamResPacket{
						Transaction_id: 0,
						Stream_id:      0,
					}
					err = (*pkt).Decode((*msg).Payload)

				case pt.RTMP_AMF0_COMMAND_RELEASE_STREAM,
					pt.RTMP_AMF0_COMMAND_FC_PUBLISH,
					pt.RTMP_AMF0_COMMAND_UNPUBLISH:
					(*pkt) = &pt.FmleStartResPacket{
						Transaction_id: 0,
					}
					err = (*pkt).Decode((*msg).Payload)

				default:
					err = fmt.Errorf("unknown request command name. %s", req_command_name)
					return
				}

				if err != nil {
					log.Println("decode result or error response msg failed.")
					return
				}

			}
		case pt.RTMP_AMF0_COMMAND_CONNECT:
			*pkt = &pt.ConnectPacket{}
			err = (*pkt).Decode((*msg).Payload)
		case pt.RTMP_AMF0_COMMAND_CREATE_STREAM:
			*pkt = &pt.CreateStreamPacket{}
			err = (*pkt).Decode((*msg).Payload)
		case pt.RTMP_AMF0_COMMAND_PLAY:
			*pkt = &pt.PlayPacket{}
			err = (*pkt).Decode((*msg).Payload)
		case pt.RTMP_AMF0_COMMAND_PAUSE:
			*pkt = &pt.PausePacket{}
			err = (*pkt).Decode((*msg).Payload)
		case pt.RTMP_AMF0_COMMAND_RELEASE_STREAM, pt.RTMP_AMF0_COMMAND_FC_PUBLISH, pt.RTMP_AMF0_COMMAND_UNPUBLISH:
			*pkt = &pt.FmleStartPacket{}
			err = (*pkt).Decode((*msg).Payload)
		case pt.RTMP_AMF0_COMMAND_PUBLISH:
			*pkt = &pt.PublishPacket{}
			err = (*pkt).Decode((*msg).Payload)
		case pt.RTMP_AMF0_DATA_SET_DATAFRAME, pt.RTMP_AMF0_DATA_ON_METADATA:
			*pkt = &pt.OnMetaDataPacket{}
			err = (*pkt).Decode((*msg).Payload)
		case pt.RTMP_AMF0_DATA_ON_CUSTOMDATA:
			(*pkt) = &pt.OnCustomDataPakcet{}
			err = (*pkt).Decode((*msg).Payload)
		case pt.RTMP_AMF0_COMMAND_CLOSE_STREAM:
			*pkt = &pt.CloseStreamPacket{}
			err = (*pkt).Decode((*msg).Payload)
		default:
			if (*msg).Header.IsAmf0Command() || (*msg).Header.IsAmf3Command() {
				*pkt = &pt.CallPacket{}
				err = (*pkt).Decode((*msg).Payload)
			}
		}

		if err != nil {
			return
		}

	} else if (*msg).Header.IsUserControlMessage() {
		*pkt = &pt.UserControlPacket{}
		err = (*pkt).Decode((*msg).Payload)
	} else if (*msg).Header.IsWindowAckledgementSize() {
		*pkt = &pt.SetWindowAckSizePacket{}
		err = (*pkt).Decode((*msg).Payload)
	} else if (*msg).Header.IsSetChunkSize() {
		*pkt = &pt.SetChunkSizePacket{}
		err = (*pkt).Decode((*msg).Payload)
	} else {
		if !(*msg).Header.IsAckledgement() && !(*msg).Header.IsSetPeerBandwidth() {
			//drop msg
		}
	}

	if err != nil {
		return
	}

	return
}
