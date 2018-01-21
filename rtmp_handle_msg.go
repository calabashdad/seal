package main

import (
	"UtilsTools/identify_panic"
	"encoding/binary"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func (rtmp *RtmpConn) HanleMsg(chunkStreamId uint32) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}
	}()

	chunk := rtmp.chunks[chunkStreamId]
	if nil == chunk {
		err = fmt.Errorf("HanleMsg:can not find the chunk strema id in chuns.")
		return
	}

	switch chunk.msg.header.typeId {
	case RTMP_MSG_AMF3CommandMessage, RTMP_MSG_AMF0CommandMessage, RTMP_MSG_AMF0DataMessage, RTMP_MSG_AMF3DataMessage:
		err = rtmp.handleAMFCommandAndDataMessage(&chunk.msg)
	case RTMP_MSG_UserControlMessage:
		err = rtmp.handleUserControlMessage(&chunk.msg)
	case RTMP_MSG_WindowAcknowledgementSize:
		err = rtmp.handleSetWindowAcknowledgementSize(&chunk.msg)
	case RTMP_MSG_SetChunkSize:
		err = rtmp.handleSetChunkSize(&chunk.msg)
	case RTMP_MSG_SetPeerBandwidth:
		err = rtmp.handleSetPeerBandWidth(&chunk.msg)
	case RTMP_MSG_Acknowledgement:
		err = rtmp.handleAcknowlegement(&chunk.msg)
	case RTMP_MSG_AbortMessage:
		err = rtmp.handleAbortMsg(&chunk.msg)
	case RTMP_MSG_EdgeAndOriginServerCommand:
		err = rtmp.handleEdgeAndOriginServerCommand(&chunk.msg)
	default:
		log.Println("unknown msg type. ", chunk.msg.header.typeId)
	}

	if err != nil {
		return
	}

	return
}

func (rtmp *RtmpConn) handleAMF0CommandResultError(msg *MessageStream) (err error) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	var offset uint32

	var transactionId float64
	err, transactionId = Amf0ReadNumber(msg.payload, &offset)
	if err != nil {
		return
	}

	requestCommand, ok := rtmp.transactionIds[transactionId]
	if !ok {
		err = fmt.Errorf("handleAMF0CommandResultError can not find the transaction id.")
		return
	}

	switch requestCommand {
	case RTMP_AMF0_COMMAND_CONNECT:
		//todo
	case RTMP_AMF0_COMMAND_CREATE_STREAM:
		//todo
	case RTMP_AMF0_COMMAND_RELEASE_STREAM, RTMP_AMF0_COMMAND_FC_PUBLISH, RTMP_AMF0_COMMAND_UNPUBLISH:
		//todo
	case RTMP_AMF0_COMMAND_ENABLEVIDEO:
		//todo
	case RTMP_AMF0_COMMAND_INSERT_KEYFREAME:
		//todo
	default:
		err = fmt.Errorf("handleAMF0CommandResultError, unknown request command name,", requestCommand)
	}

	if err != nil {
		return
	}

	return
}

func (rtmp *RtmpConn) handleAMFCommandAndDataMessage(msg *MessageStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}
	}()

	log.Println("amf0/3 command or amf0/3 data")

	var offset uint32

	// skip 1bytes to decode the amf3 command.
	if RTMP_MSG_AMF3CommandMessage == msg.header.typeId && msg.header.length >= 1 {
		offset += 1
	}

	//read the command name.
	err, command := Amf0ReadString(msg.payload, &offset)
	if err != nil {
		return
	}

	if 0 == len(command) {
		err = fmt.Errorf("Amf0ReadString failed, command is nil.")
		return
	}

	log.Println("msg typeid=", msg.header.typeId, ",command =", command)

	switch command {
	case RTMP_AMF0_COMMAND_RESULT, RTMP_AMF0_COMMAND_ERROR:
		err = rtmp.handleAMF0CommandResultError(msg)
	case RTMP_AMF0_COMMAND_CONNECT:
		err = rtmp.handleAMF0CommandConnect(msg)
	case RTMP_AMF0_COMMAND_CREATE_STREAM:
		//todo.
	case RTMP_AMF0_COMMAND_PLAY:
		//todo.
	case RTMP_AMF0_COMMAND_PAUSE:
		//todo.
	case RTMP_AMF0_COMMAND_RELEASE_STREAM:
		//todo.
	case RTMP_AMF0_COMMAND_FC_PUBLISH:
		//todo.
	case RTMP_AMF0_COMMAND_PUBLISH:
		//todo.
	case RTMP_AMF0_COMMAND_UNPUBLISH:
		//todo.
	case RTMP_AMF0_COMMAND_KEEPLIVE:
		//todo.
	case RTMP_AMF0_COMMAND_ENABLEVIDEO:
		//todo.
	case RTMP_AMF0_DATA_SET_DATAFRAME, RTMP_AMF0_DATA_ON_METADATA:
		//todo.
	case RTMP_AMF0_DATA_ON_CUSTOMDATA:
		//todo.
	case RTMP_AMF0_COMMAND_CLOSE_STREAM:
		//todo.
	case RTMP_AMF0_COMMAND_ON_BW_DONE:
		//todo
	case RTMP_AMF0_COMMAND_ON_STATUS:
		//todo
	case RTMP_AMF0_COMMAND_INSERT_KEYFREAME:
		//todo
	case RTMP_AMF0_DATA_SAMPLE_ACCESS:
		//todo.
	default:
		log.Println("handleAMFCommandAndDataMessage:unknown command name.", command)
	}

	if err != nil {
		return
	}

	return
}

type Amf0CommandConnectPkg struct {
	command        string
	transactionId  float64
	commandObjects []Amf0Object
	amfOptional    interface{}
}

func (rtmp *RtmpConn) handleAMF0CommandConnect(msg *MessageStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	var connectPkg Amf0CommandConnectPkg

	var offset uint32

	err, connectPkg.command = Amf0ReadString(msg.payload, &offset)
	if err != nil {
		return
	}

	if connectPkg.command != RTMP_AMF0_COMMAND_CONNECT {
		err = fmt.Errorf("handleAMF0CommandConnect command is error. command=", connectPkg.command)
		return
	}

	err, connectPkg.transactionId = Amf0ReadNumber(msg.payload, &offset)
	if err != nil {
		return
	}

	//this method is not strict for float type. just a warn.
	if 1 != connectPkg.transactionId {
		log.Println("warn:handleAMF0CommandConnect: transactionId is not 1. transactionId=", connectPkg.transactionId)
	}

	err, connectPkg.commandObjects = Amf0ReadObject(msg.payload, &offset)
	if err != nil {
		return
	}

	if offset < uint32(len(msg.payload)) {
		var v interface{}
		var marker uint8
		err, v, marker = Amf0Discovery(msg.payload, &offset)
		if err != nil {
			return
		}

		if RTMP_AMF0_Object == marker {
			connectPkg.amfOptional = v
		}
	}

	err = rtmp.ParseConnectPkg(&connectPkg)
	if err != nil {
		log.Println("parse connect pkg error.", err)
		return
	}

	err = rtmp.CommonMsgSetWindowAcknowledgementSize(msg.header.preferCsId, 2500000)
	if err != nil {
		return
	}

	err = rtmp.CommonMsgSetPeerBandwidth(msg.header.preferCsId, 2500000, 2)
	if err != nil {
		return
	}

	err = rtmp.ResponseConnectApp(msg.header.preferCsId)
	if err != nil {
		return
	}

	return
}

func (pkg *Amf0CommandConnectPkg) Amf0ObjectsGetProperty(key string) (value interface{}) {

	for _, v := range pkg.commandObjects {
		if v.propertyName == key {
			return v.value
		}
	}

	return
}

func (rtmp *RtmpConn) ParseConnectPkg(pkg *Amf0CommandConnectPkg) (err error) {
	tcUrlValue := pkg.Amf0ObjectsGetProperty("tcUrl")
	if nil == tcUrlValue {
		err = fmt.Errorf("tcUrl is nil.")
		return
	}

	objectEncodingValue := pkg.Amf0ObjectsGetProperty("objectEncoding")
	if objectEncodingValue != nil {
		rtmp.objectEncoding = objectEncodingValue.(float64)
	}

	var rtmpUrlData RtmpUrlData
	err = rtmpUrlData.ParseUrl(tcUrlValue.(string))
	if err != nil {
		return
	}

	err = rtmpUrlData.Discover()
	if err != nil {
		return
	}

	return
}

type RtmpUrlData struct {
	schema string
	host   string
	port   uint16
	app    string
	stream string
	token  string
}

//format: rtmp://127.0.0.1:1935/live/test?token=abc123
func (urlData *RtmpUrlData) ParseUrl(url string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}
	}()

	//url is not contain the stream with token.
	var urlTmp string

	urlTmp = strings.Replace(url, "://", " ", 1)
	urlTmp = strings.Replace(urlTmp, ":", " ", 1)
	urlTmp = strings.Replace(urlTmp, "/", " ", 1)

	urlSplit := strings.Split(urlTmp, " ")

	if 4 == len(urlSplit) {
		urlData.schema = urlSplit[0]
		urlData.host = urlSplit[1]
		port, ok := strconv.Atoi(urlSplit[2])
		if nil == ok {
			//the port is not default
			if port > 0 && port < 65536 {
				urlData.port = uint16(port)
			} else {
				err = fmt.Errorf("tcUrl port format is error, port=", port)
				return
			}

		} else {
			err = fmt.Errorf("tcurl format error when convert port format, err=", ok)
			return
		}
		urlData.app = urlSplit[3]
	} else {
		err = fmt.Errorf("tcUrl format is error. tcUrl=", url)
		return
	}

	return
}

func (urlData *RtmpUrlData) Discover() (err error) {
	if 0 == len(urlData.schema) ||
		0 == len(urlData.host) ||
		0 == len(urlData.app) {
		err = fmt.Errorf("discover url data failed. url data=", urlData)
		return
	}
	return
}

func (rtmp *RtmpConn) CommonMsgSetWindowAcknowledgementSize(chunkStreamId uint32, WindowAcknowledgementSize uint32) (err error) {

	var msg MessageStream

	//msg payload
	msg.payload = make([]uint8, 4)
	binary.BigEndian.PutUint32(msg.payload[:], WindowAcknowledgementSize)

	//msg header
	msg.header.length = 4
	msg.header.typeId = RTMP_MSG_WindowAcknowledgementSize
	msg.header.streamId = 0
	msg.header.preferCsId = chunkStreamId

	//begin to send.
	if msg.header.preferCsId < 2 {
		msg.header.preferCsId = RTMP_CID_ProtocolControl
	}

	err = rtmp.SendMsg(&msg)
	if err != nil {
		return
	}

	return
}

func (rtmp *RtmpConn) CommonMsgResponseWindowAcknowledgement(chunkStreamId uint32, ackedSize uint32) (err error) {

	var msg MessageStream

	//msg payload
	msg.payload = make([]uint8, 4)
	binary.BigEndian.PutUint32(msg.payload[:], ackedSize)

	//msg header
	msg.header.length = 4
	msg.header.typeId = RTMP_MSG_Acknowledgement
	msg.header.streamId = 0
	msg.header.preferCsId = chunkStreamId

	//begin to send.
	if msg.header.preferCsId < 2 {
		msg.header.preferCsId = RTMP_CID_ProtocolControl
	}

	err = rtmp.SendMsg(&msg)
	if err != nil {
		return
	}

	return
}

func (rtmp *RtmpConn) CommonMsgSetPeerBandwidth(chunkStreamId uint32, bandWidthValue uint32, limitType uint8) (err error) {

	var msg MessageStream

	//msg payload
	msg.payload = make([]uint8, 5)
	binary.BigEndian.PutUint32(msg.payload[:4], bandWidthValue)
	msg.payload[4] = limitType

	//msg header
	msg.header.length = 4
	msg.header.typeId = RTMP_MSG_SetPeerBandwidth
	msg.header.streamId = 0
	msg.header.preferCsId = chunkStreamId

	//begin to send.
	if msg.header.preferCsId < 2 {
		msg.header.preferCsId = RTMP_CID_ProtocolControl
	}

	err = rtmp.SendMsg(&msg)
	if err != nil {
		return
	}

	return
}

func (rtmp *RtmpConn) CommonMsgSetChunkSize(chunkSize uint32) (err error) {
	var msg MessageStream

	//msg payload
	msg.payload = make([]uint8, 4)
	binary.BigEndian.PutUint32(msg.payload[:], chunkSize)

	//msg header
	msg.header.length = 4
	msg.header.typeId = RTMP_MSG_SetChunkSize
	msg.header.streamId = 0
	msg.header.preferCsId = RTMP_CID_ProtocolControl

	err = rtmp.SendMsg(&msg)
	if err != nil {
		return
	}

	return
}

/**
* 5.4. User Control Message (4)
*
* for the EventData is 4bytes.
* Stream Begin(=0)              4-bytes stream ID
* Stream EOF(=1)                4-bytes stream ID
* StreamDry(=2)                 4-bytes stream ID
* SetBufferLength(=3)           8-bytes 4bytes stream ID, 4bytes buffer length.
* StreamIsRecorded(=4)          4-bytes stream ID
* PingRequest(=6)               4-bytes timestamp local server time
* PingResponse(=7)              4-bytes timestamp received ping request.
*
* 3.7. User Control message
* +------------------------------+-------------------------
* | Event Type ( 2- bytes ) | Event Data
* +------------------------------+-------------------------
* Figure 5 Pay load for the ‘User Control Message’.
 */
type UserControlMsg struct {
	eventType uint16
	eventData uint32
	/**
	 * 4bytes if event_type is SetBufferLength; otherwise 0.
	 */
	extraData uint32
}

func (rtmp *RtmpConn) handleUserControlMessage(msg *MessageStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	log.Println("msg type usercontrol msg.")

	var offset uint32

	if uint32(len(msg.payload))-offset < (2 + 4) {
		err = fmt.Errorf("handleUserControlMessage, 0, length is not enough.")
		return
	}

	var userCtrlMsg UserControlMsg

	userCtrlMsg.eventType = binary.BigEndian.Uint16(msg.payload[offset : offset+2])
	offset += 2

	userCtrlMsg.eventData = binary.BigEndian.Uint32(msg.payload[offset : offset+4])
	offset += 4

	log.Println("user control msg, type=", userCtrlMsg.eventType)

	if SrcPCUCSetBufferLength == userCtrlMsg.eventType {
		if uint32(len(msg.payload))-offset < 4 {
			err = fmt.Errorf("handleUserControlMessage, 1, length is not enough.")
			return
		}

		userCtrlMsg.extraData = binary.BigEndian.Uint32(msg.payload[offset : offset+4])
		offset += 4
	}

	if err != nil {
		return
	}

	switch userCtrlMsg.eventType {
	case SrcPCUCStreamBegin:
	case SrcPCUCStreamEOF:
	case SrcPCUCStreamDry:
	case SrcPCUCSetBufferLength:
	case SrcPCUCStreamIsRecorded:
	case SrcPCUCPingRequest:
		err = rtmp.ResponsePingMessage(msg.header.preferCsId, &userCtrlMsg)
	case SrcPCUCPingResponse:
	default:
		log.Println("HandleMsgUserControl unknown event type.type=", userCtrlMsg.eventType)
	}

	if err != nil {
		return
	}

	return
}

func (rtmp *RtmpConn) ResponsePingMessage(chunkStreamId uint32, userCtrl *UserControlMsg) (err error) {
	var msg MessageStream

	//msg payload
	var offset uint32

	msg.payload = make([]uint8, 2+4) // 2(type) + 4(data)
	binary.BigEndian.PutUint16(msg.payload[offset:offset+2], SrcPCUCPingResponse)
	offset += 2
	binary.BigEndian.PutUint32(msg.payload[offset:offset+4], userCtrl.eventData)
	offset += 4

	//msg header
	msg.header.length = uint32(len(msg.payload))
	msg.header.typeId = RTMP_MSG_UserControlMessage
	msg.header.streamId = 0
	if chunkStreamId < 2 {
		msg.header.preferCsId = RTMP_CID_ProtocolControl
	} else {
		msg.header.preferCsId = chunkStreamId
	}

	err = rtmp.SendMsg(&msg)
	if err != nil {
		return
	}

	return
}

func (rtmp *RtmpConn) handleSetWindowAcknowledgementSize(msg *MessageStream) (err error) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	var windowAcknowlegementSize uint32
	if len(msg.payload) >= 4 {
		windowAcknowlegementSize = binary.BigEndian.Uint32(msg.payload[0:4])
	} else {
		err = fmt.Errorf("handleSetWindowAcknowledgementSize payload len < 4", len(msg.payload))
		return
	}

	if windowAcknowlegementSize > 0 {
		rtmp.ackWindow.ackWindowSize = windowAcknowlegementSize
		log.Println("peer set window acknowlegement size.", windowAcknowlegementSize)

	} else {
		//ignored.
		log.Println("HandleMsgSetWindowsAcknowlegementSize, ack size is invalied.", windowAcknowlegementSize)
	}

	return
}

func (rtmp *RtmpConn) handleSetChunkSize(msg *MessageStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	var chunkSize uint32

	if len(msg.payload) >= 4 {
		chunkSize = binary.BigEndian.Uint32(msg.payload[0:4])
	} else {
		err = fmt.Errorf("handleSetChunkSize payload length < 4", len(msg.payload))
		return
	}

	if chunkSize >= RTMP_CHUNKSIZE_MIN && chunkSize <= RTMP_CHUNKSIZE_MAX {
		rtmp.chunkSize = chunkSize
		log.Println("peer set chunk size success. chunk size=", chunkSize)
	} else {
		//ignored
		log.Println("HandleMsgSetChunkSize, chunk size is invalid.", chunkSize)
	}

	return
}

func (rtmp *RtmpConn) handleSetPeerBandWidth(msg *MessageStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	log.Println("set bandwidth")

	return
}

func (rtmp *RtmpConn) handleAcknowlegement(msg *MessageStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	log.Println("acknowlegement")

	return
}

func (rtmp *RtmpConn) handleAbortMsg(msg *MessageStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	log.Println("abort msg, remote=", rtmp.Conn.RemoteAddr())

	return
}

func (rtmp *RtmpConn) handleEdgeAndOriginServerCommand(msg *MessageStream) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at,", identify_panic.IdentifyPanic())
		}
	}()

	log.Println("edge and origin server command, remote=", rtmp.Conn.RemoteAddr())

	return
}
