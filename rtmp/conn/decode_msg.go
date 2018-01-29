package conn

import (
	"UtilsTools/identify_panic"
	"fmt"
	"log"
	"seal/rtmp/protocol"
)

func (rc *RtmpConn) DecodeMsg(msg **protocol.Message, pkt *protocol.Packet) (err error) {
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

	if (*msg).Header.Is_amf0_command() || (*msg).Header.Is_amf3_command() ||
		(*msg).Header.Is_amf0_data() || (*msg).Header.Is_amf3_data() {

		if (*msg).Header.Is_amf3_command() && len((*msg).Payload) > 1 {
			offset += 1
		}

		var command string
		err, command = protocol.Amf0ReadString((*msg).Payload, &offset)
		if err != nil {
			log.Println("read command failed when decode msg.")
			return
		}

		switch command {
		//command: result or error.
		case protocol.RTMP_AMF0_COMMAND_RESULT, protocol.RTMP_AMF0_COMMAND_ERROR:
			if true {
				var transaction_id float64
				err, transaction_id = protocol.Amf0ReadNumber((*msg).Payload, &offset)
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
				case protocol.RTMP_AMF0_COMMAND_CONNECT:
					*pkt = &protocol.ConnectResPacket{}
					err = (*pkt).Decode((*msg).Payload)

				case protocol.RTMP_AMF0_COMMAND_CREATE_STREAM:
					(*pkt) = &protocol.CreateStreamResPacket{
						Transaction_id: 0,
						Stream_id:      0,
					}
					err = (*pkt).Decode((*msg).Payload)

				case protocol.RTMP_AMF0_COMMAND_RELEASE_STREAM,
					protocol.RTMP_AMF0_COMMAND_FC_PUBLISH,
					protocol.RTMP_AMF0_COMMAND_UNPUBLISH:
					(*pkt) = &protocol.FmleStartResPacket{
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
		case protocol.RTMP_AMF0_COMMAND_CONNECT:
			*pkt = &protocol.ConnectPacket{}
			err = (*pkt).Decode((*msg).Payload)
		case protocol.RTMP_AMF0_COMMAND_CREATE_STREAM:
			*pkt = &protocol.CreateStreamPacket{}
			err = (*pkt).Decode((*msg).Payload)
		case protocol.RTMP_AMF0_COMMAND_PLAY:
			*pkt = &protocol.PlayPacket{}
			err = (*pkt).Decode((*msg).Payload)
		case protocol.RTMP_AMF0_COMMAND_PAUSE:
			*pkt = &protocol.PausePacket{}
			err = (*pkt).Decode((*msg).Payload)
		case protocol.RTMP_AMF0_COMMAND_RELEASE_STREAM, protocol.RTMP_AMF0_COMMAND_FC_PUBLISH, protocol.RTMP_AMF0_COMMAND_UNPUBLISH:
			*pkt = &protocol.FmleStartPacket{}
			err = (*pkt).Decode((*msg).Payload)
		case protocol.RTMP_AMF0_COMMAND_PUBLISH:
			*pkt = &protocol.PublishPacket{}
			err = (*pkt).Decode((*msg).Payload)
		case protocol.RTMP_AMF0_DATA_SET_DATAFRAME, protocol.RTMP_AMF0_DATA_ON_METADATA:
			*pkt = &protocol.OnMetaDataPacket{}
			err = (*pkt).Decode((*msg).Payload)
		case protocol.RTMP_AMF0_DATA_ON_CUSTOMDATA:
			(*pkt) = &protocol.OnCustomDataPakcet{}
			err = (*pkt).Decode((*msg).Payload)
		case protocol.RTMP_AMF0_COMMAND_CLOSE_STREAM:
			*pkt = &protocol.CloseStreamPacket{}
			err = (*pkt).Decode((*msg).Payload)
		default:
			if (*msg).Header.Is_amf0_command() || (*msg).Header.Is_amf3_command() {
				*pkt = &protocol.CallPacket{}
				err = (*pkt).Decode((*msg).Payload)
			}
		}

		if err != nil {
			return
		}

	} else if (*msg).Header.Is_user_control_message() {
		*pkt = &protocol.UserControlPacket{}
		err = (*pkt).Decode((*msg).Payload)
	} else if (*msg).Header.Is_window_ackledgement_size() {
		*pkt = &protocol.SetWindowAckSizePacket{}
		err = (*pkt).Decode((*msg).Payload)
	} else if (*msg).Header.Is_set_chunk_size() {
		*pkt = &protocol.SetChunkSizePacket{}
		err = (*pkt).Decode((*msg).Payload)
	} else {
		if !(*msg).Header.Is_ackledgement() && !(*msg).Header.Is_set_peer_bandwidth() {
			//drop msg
		}
	}

	if err != nil {
		return
	}

	return
}
