package protocol_stack

const (
	/**
	 * 6.1.2. Chunk MessageStream Header
	 * There are four different formats for the chunk message header,
	 * selected by the "chunkFmt" field in the chunk basic header.
	 */
	// 6.1.2.1. Type 0
	// Chunks of Type 0 are 11 bytes long. This type MUST be used at the
	// start of a chunk stream, and whenever the stream timestampDelta goes
	// backward (e.g., because of a backward seek).
	RTMP_FMT_TYPE0 = 0
	// 6.1.2.2. Type 1
	// Chunks of Type 1 are 7 bytes long. The message stream ID is not
	// included; this chunk takes the same stream ID as the preceding chunk.
	// Streams with variable-sized messages (for example, many video
	// formats) SHOULD use this format for the first chunk of each new
	// message after the first.
	RTMP_FMT_TYPE1 = 1
	// 6.1.2.3. Type 2
	// Chunks of Type 2 are 3 bytes long. Neither the stream ID nor the
	// message length is included; this chunk has the same stream ID and
	// message length as the preceding chunk. Streams with constant-sized
	// messages (for example, some audio and data formats) SHOULD use this
	// format for the first chunk of each message after the first.
	RTMP_FMT_TYPE2 = 2
	// 6.1.2.4. Type 3
	// Chunks of Type 3 have no header. Stream ID, message length and
	// timestampDelta delta are not present; chunks of this type take values from
	// the preceding chunk. When a single message is split into chunks, all
	// chunks of a message except the first one, SHOULD use this type. Refer
	// to example 2 in section 6.2.2. Stream consisting of messages of
	// exactly the same size, stream ID and spacing in time SHOULD use this
	// type for all chunks after chunk of Type 2. Refer to example 1 in
	// section 6.2.1. If the delta between the first message and the second
	// message is same as the time stamp of first message, then chunk of
	// type 3 would immediately follow the chunk of type 0 as there is no
	// need for a chunk of type 2 to register the delta. If Type 3 chunk
	// follows a Type 0 chunk, then timestampDelta delta for this Type 3 chunk is
	// the same as the timestampDelta of Type 0 chunk.
	RTMP_FMT_TYPE3 = 3
)

const (
	/**
	 * the chunk stream id used for some under-layer message,
	 * for example, the PC(protocol control) message.
	 */
	RTMP_CID_ProtocolControl = 0x02
	/**
	 * the AMF0/AMF3 command message, invoke method and return the result, over NetConnection.
	 * generally use 0x03.
	 */
	RTMP_CID_OverConnection = 0x03
	/**
	 * the AMF0/AMF3 command message, invoke method and return the result, over NetConnection,
	 * the midst state(we guess).
	 * rarely used, e.g. onStatus(NetStream.Play.Reset).
	 */
	RTMP_CID_OverConnection2 = 0x04
	/**
	 * the stream message(amf0/amf3), over NetStream.
	 * generally use 0x05.
	 */
	RTMP_CID_OverStream = 0x05
	/**
	 * the stream message(amf0/amf3), over NetStream, the midst state(we guess).
	 * rarely used, e.g. play("mp4:mystram.f4v")
	 */
	RTMP_CID_OverStream2 = 0x08
	/**
	 * the stream message(video), over NetStream
	 * generally use 0x06.
	 */
	RTMP_CID_Video = 0x06
	/**
	 * the stream message(audio), over NetStream.
	 * generally use 0x07.
	 */
	RTMP_CID_Audio = 0x07
)

const (
	RTMP_EXTENDED_TIMESTAMP = 0xFFFFFF
)

const (
	RTMP_DEFAULT_CHUNK_SIZE = 128

	RTMP_CHUNKSIZE_MIN = 128
	RTMP_CHUNKSIZE_MAX = 65536
)

const (
	/**
	3. Types of messages
	The server and the client send messages over the network to
	communicate with each other. The messages can be of any type which
	includes audio messages, video messages, command messages, shared
	object messages, data messages, and user control messages.
	3.1. Command message
	Command messages carry the AMF-encoded commands between the client
	and the server. These messages have been assigned message type value
	of 20 for AMF0 encoding and message type value of 17 for AMF3
	encoding. These messages are sent to perform some operations like
	connect, createStream, publish, play, pause on the peer. Command
	messages like onstatus, result etc. are used to inform the sender
	about the status of the requested commands. A command message
	consists of command name, transaction ID, and command object that
	contains related parameters. A client or a server can request Remote
	Procedure Calls (RPC) over streams that are communicated using the
	command messages to the peer.
	*/

	RTMP_MSG_AMF3CommandMessage = 17 // 0x11
	RTMP_MSG_AMF0CommandMessage = 20 // 0x14

	/**
	3.2. Data message
	The client or the server sends this message to send Metadata or any
	user data to the peer. Metadata includes details about the
	data(audio, video etc.) like creation time, duration, theme and so
	on. These messages have been assigned message type value of 18 for
	AMF0 and message type value of 15 for AMF3.
	*/

	RTMP_MSG_AMF0DataMessage = 18 // 0x12
	RTMP_MSG_AMF3DataMessage = 15 // 0x0F
)

const (
	/**
	5. Protocol Control Messages
	RTMP reserves message type IDs 1-7 for protocol control messages.
	These messages contain information needed by the RTM Chunk Stream
	protocol or RTMP itself. Protocol messages with IDs 1 & 2 are
	reserved for usage with RTM Chunk Stream protocol. Protocol messages
	with IDs 3-6 are reserved for usage of RTMP. Protocol message with ID
	7 is used between edge server and origin server.
	*/
	RTMP_MSG_SetChunkSize               = 0x01
	RTMP_MSG_AbortMessage               = 0x02
	RTMP_MSG_Acknowledgement            = 0x03
	RTMP_MSG_UserControlMessage         = 0x04
	RTMP_MSG_WindowAcknowledgementSize  = 0x05
	RTMP_MSG_SetPeerBandwidth           = 0x06
	RTMP_MSG_EdgeAndOriginServerCommand = 0x07
)

const (
	// AMF0 marker
	RTMP_AMF0_Number      = 0x00
	RTMP_AMF0_Boolean     = 0x01
	RTMP_AMF0_String      = 0x02
	RTMP_AMF0_Object      = 0x03
	RTMP_AMF0_MovieClip   = 0x04 // reserved, not supported
	RTMP_AMF0_Null        = 0x05
	RTMP_AMF0_Undefined   = 0x06
	RTMP_AMF0_Reference   = 0x07
	RTMP_AMF0_EcmaArray   = 0x08
	RTMP_AMF0_ObjectEnd   = 0x09
	RTMP_AMF0_StrictArray = 0x0A
	RTMP_AMF0_Date        = 0x0B
	RTMP_AMF0_LongString  = 0x0C
	RTMP_AMF0_UnSupported = 0x0D
	RTMP_AMF0_RecordSet   = 0x0E // reserved, not supported
	RTMP_AMF0_XmlDocument = 0x0F
	RTMP_AMF0_TypedObject = 0x10
	// AVM+ object is the AMF3 object.
	RTMP_AMF0_AVMplusObject = 0x11
	// origin array whos data takes the same form as LengthValueBytes
	RTMP_AMF0_OriginStrictArray = 0x20

	// User defined
	RTMP_AMF0_Invalid = 0x3F
)

const (
	/**
	 * amf0 command message, command name macros
	 */
	RTMP_AMF0_COMMAND_CONNECT          = "connect"
	RTMP_AMF0_COMMAND_CREATE_STREAM    = "createStream"
	RTMP_AMF0_COMMAND_CLOSE_STREAM     = "closeStream"
	RTMP_AMF0_COMMAND_PLAY             = "play"
	RTMP_AMF0_COMMAND_PAUSE            = "pause"
	RTMP_AMF0_COMMAND_ON_BW_DONE       = "onBWDone"
	RTMP_AMF0_COMMAND_ON_STATUS        = "onStatus"
	RTMP_AMF0_COMMAND_RESULT           = "_result"
	RTMP_AMF0_COMMAND_ERROR            = "_error"
	RTMP_AMF0_COMMAND_RELEASE_STREAM   = "releaseStream"
	RTMP_AMF0_COMMAND_FC_PUBLISH       = "FCPublish"
	RTMP_AMF0_COMMAND_UNPUBLISH        = "FCUnpublish"
	RTMP_AMF0_COMMAND_PUBLISH          = "publish"
	RTMP_AMF0_COMMAND_KEEPLIVE         = "JMS.KeepAlive"
	RTMP_AMF0_COMMAND_ENABLEVIDEO      = "JMS.EnableVideo"
	RTMP_AMF0_COMMAND_INSERT_KEYFREAME = "JMS.InsertKeyframe"
	RTMP_AMF0_DATA_SAMPLE_ACCESS       = "|RtmpSampleAccess"
	RTMP_AMF0_DATA_SET_DATAFRAME       = "@setDataFrame"
	RTMP_AMF0_DATA_ON_METADATA         = "onMetaData"
	RTMP_AMF0_DATA_ON_CUSTOMDATA       = "onCustomData"
)

const (
	/**
	 * band width check method name, which will be invoked by client.
	 * band width check mothods use SrsBandwidthPacket as its internal packet type,
	 * so ensure you set command name when you use it.
	 */
	// server play control
	SRS_BW_CHECK_START_PLAY    = "onSrsBandCheckStartPlayBytes"
	SRS_BW_CHECK_STARTING_PLAY = "onSrsBandCheckStartingPlayBytes"
	SRS_BW_CHECK_STOP_PLAY     = "onSrsBandCheckStopPlayBytes"
	SRS_BW_CHECK_STOPPED_PLAY  = "onSrsBandCheckStoppedPlayBytes"

	// server publish control
	SRS_BW_CHECK_START_PUBLISH    = "onSrsBandCheckStartPublishBytes"
	SRS_BW_CHECK_STARTING_PUBLISH = "onSrsBandCheckStartingPublishBytes"
	SRS_BW_CHECK_STOP_PUBLISH     = "onSrsBandCheckStopPublishBytes"
	// @remark, flash never send out this packet, for its queue is full.
	SRS_BW_CHECK_STOPPED_PUBLISH = "onSrsBandCheckStoppedPublishBytes"

	// EOF control.
	// the report packet when check finished.
	SRS_BW_CHECK_FINISHED = "onSrsBandCheckFinished"
	// @remark, flash never send out this packet, for its queue is full.
	SRS_BW_CHECK_FINAL = "finalClientPacket"

	// data packets
	SRS_BW_CHECK_PLAYING    = "onSrsBandCheckPlaying"
	SRS_BW_CHECK_PUBLISHING = "onSrsBandCheckPublishing"
)

const (
	/**
	 * onStatus consts.
	 */
	StatusLevel       = "level"
	StatusCode        = "code"
	StatusDescription = "description"
	StatusDetails     = "details"
	StatusClientId    = "clientid"

	// status value
	StatusLevelStatus = "status"

	// status error
	StatusLevelError = "error"

	// code value
	StatusCodeConnectSuccess   = "NetConnection.Connect.Success"
	StatusCodeConnectRejected  = "NetConnection.Connect.Rejected"
	StatusCodeStreamReset      = "NetStream.Play.Reset"
	StatusCodeStreamStart      = "NetStream.Play.Start"
	StatusCodeStreamPause      = "NetStream.Pause.Notify"
	StatusCodeStreamUnpause    = "NetStream.Unpause.Notify"
	StatusCodePublishStart     = "NetStream.Publish.Start"
	StatusCodeDataStart        = "NetStream.Data.Start"
	StatusCodeUnpublishSuccess = "NetStream.Unpublish.Success"

	// FMLE
	RTMP_AMF0_COMMAND_ON_FC_PUBLISH   = "onFCPublish"
	RTMP_AMF0_COMMAND_ON_FC_UNPUBLISH = "onFCUnpublish"

	// default stream id for response the createStream request.
	SRS_DEFAULT_SID = 1
)

const (
	//objectEncoding default value.
	RTMP_SIG_AMF0_VER = 0.0
)

const (
	/**
	 * The server sends this event to notify the client
	 * that a stream has become functional and can be
	 * used for communication. By default, this event
	 * is sent on ID 0 after the application connect
	 * command is successfully received from the
	 * client. The event data is 4-byte and represents
	 * the stream ID of the stream that became
	 * functional.
	 */
	SrcPCUCStreamBegin = 0x00

	/**
	* The server sends this event to notify the client
	* that the playback of data is over as requested
	* on this stream. No more data is sent without
	* issuing additional commands. The client discards
	* the messages received for the stream. The
	* 4 bytes of event data represent the ID of the
	* stream on which playback has ended.
	 */
	SrcPCUCStreamEOF = 0x01

	/**
	 * The server sends this event to notify the client
	 * that there is no more data on the stream. If the
	 * server does not detect any message for a time
	 * period, it can notify the subscribed clients
	 * that the stream is dry. The 4 bytes of event
	 * data represent the stream ID of the dry stream.
	 */
	SrcPCUCStreamDry = 0x02

	/**
	 * The client sends this event to inform the server
	 * of the buffer size (in milliseconds) that is
	 * used to buffer any data coming over a stream.
	 * This event is sent before the server starts
	 * processing the stream. The first 4 bytes of the
	 * event data represent the stream ID and the next
	 * 4 bytes represent the buffer length, in
	 * milliseconds.
	 */
	SrcPCUCSetBufferLength = 0x03 // 8bytes event-data

	/**
	 * The server sends this event to notify the client
	 * that the stream is a recorded stream. The
	 * 4 bytes event data represent the stream ID of
	 * the recorded stream.
	 */
	SrcPCUCStreamIsRecorded = 0x04

	/**
	 * The server sends this event to test whether the
	 * client is reachable. Event data is a 4-byte
	 * timestamp, representing the local server time
	 * when the server dispatched the command. The
	 * client responds with kMsgPingResponse on
	 * receiving kMsgPingRequest.
	 */
	SrcPCUCPingRequest = 0x06

	/**
	 * The client sends this event to the server in
	 * response to the ping request. The event data is
	 * a 4-byte timestamp, which was received with the
	 * kMsgPingRequest request.
	 */
	SrcPCUCPingResponse = 0x07
)

const (
	FMS_VERSION = "1.0.0"
)
