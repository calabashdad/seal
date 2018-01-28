package conn

import (
	"UtilsTools/identify_panic"
	"fmt"
	"log"
	"seal/rtmp/protocol"
)

func (rc *RtmpConn) DecodeMsg(msg *protocol.Message, pkt protocol.Packet) (err error) {
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

	if msg.Header.Is_amf0_command() || msg.Header.Is_amf3_command() ||
		msg.Header.Is_amf0_data() || msg.Header.Is_amf3_data() {

		if msg.Header.Is_amf3_command() && len(msg.Payload) > 1 {
			offset += 1
		}

		var command string
		err, command = protocol.Amf0ReadString(msg.Payload, &offset)
		if err != nil {
			log.Println("read command failed when decode msg.")
			return
		}

		if protocol.RTMP_AMF0_COMMAND_RESULT == command || protocol.RTMP_AMF0_COMMAND_ERROR == command {

			var transaction_id float64
			err, transaction_id = protocol.Amf0ReadNumber(msg.Payload, &offset)
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
				pkt = &protocol.ConnectResPacket{}
				err = pkt.Decode(msg.Payload)

			case protocol.RTMP_AMF0_COMMAND_CREATE_STREAM:
				pkt = &protocol.CreateStreamResPacket{
					Transaction_id: 0,
					Stream_id:      0,
				}
				err = pkt.Decode(msg.Payload)

			case protocol.RTMP_AMF0_COMMAND_RELEASE_STREAM,
				protocol.RTMP_AMF0_COMMAND_FC_PUBLISH,
				protocol.RTMP_AMF0_COMMAND_UNPUBLISH:
				pkt = &protocol.FmleStartResPacket{
					Transaction_id: 0,
				}
				err = pkt.Decode(msg.Payload)

			default:
				err = fmt.Errorf("unknown request command name.", req_command_name)
				return
			}

			if err != nil {
				log.Println("decode result or error response msg failed.")
				return
			}
		}
	}

	return
}
