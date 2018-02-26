package pt

const (
	// RtmpFmtType0 Chunk MessageStream Header
	// There are four different formats for the chunk message header,
	// selected by the "chunkFmt" field in the chunk basic header.
	// 6.1.2.1. Type 0
	// Chunks of Type 0 are 11 bytes long. This type MUST be used at the
	// start of a chunk stream, and whenever the stream timestampDelta goes
	// backward (e.g., because of a backward seek).
	RtmpFmtType0 = 0

	// RtmpFmtType1 Type 1
	// Chunks of Type 1 are 7 bytes long. The message stream ID is not
	// included; this chunk takes the same stream ID as the preceding chunk.
	// Streams with variable-sized messages (for example, many video
	// formats) SHOULD use this format for the first chunk of each new
	// message after the first.
	RtmpFmtType1 = 1

	//RtmpFmtType2 Type 2
	// Chunks of Type 2 are 3 bytes long. Neither the stream ID nor the
	// message length is included; this chunk has the same stream ID and
	// message length as the preceding chunk. Streams with constant-sized
	// messages (for example, some audio and data formats) SHOULD use this
	// format for the first chunk of each message after the first.
	RtmpFmtType2 = 2

	// RtmpFmtType3 Type 3
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
	RtmpFmtType3 = 3
)

const (
	// RtmpCidProtocolControl the chunk stream id used for some under-layer message,
	// for example, the PC(protocol control) message.
	RtmpCidProtocolControl = 0x02

	// RtmpCidOverConnection the AMF0/AMF3 command message, invoke method and return the result, over NetConnection.
	// generally use 0x03.
	RtmpCidOverConnection = 0x03

	// RtmpCidOverConnection2 the AMF0/AMF3 command message, invoke method and return the result, over NetConnection,
	// the midst state(we guess).
	// rarely used, e.g. onStatus(NetStream.Play.Reset).
	RtmpCidOverConnection2 = 0x04

	// RtmpCidOverStream the stream message(amf0/amf3), over NetStream.
	// generally use 0x05.
	RtmpCidOverStream = 0x05

	// RtmpCidOverStream2 the stream message(amf0/amf3), over NetStream, the midst state(we guess).
	// rarely used, e.g. play("mp4:mystram.f4v")
	RtmpCidOverStream2 = 0x08

	// RtmpCidVideo the stream message(video), over NetStream
	// generally use 0x06.
	RtmpCidVideo = 0x06

	// RtmpCidAudio  the stream message(audio), over NetStream.
	// generally use 0x07.
	RtmpCidAudio = 0x07
)

const (
	// RtmpExtendTimeStamp rtmp extend timestamp in message, when time > 0xffffff
	RtmpExtendTimeStamp = 0xFFFFFF
)

const (
	// RtmpDefalutChunkSize rmtp chunk default size
	RtmpDefalutChunkSize = 128

	// RtmpChunkSizeMin rtmp chunk min size
	RtmpChunkSizeMin = 128

	// RtmpChunkSizeMax rtmp chunk max size
	RtmpChunkSizeMax = 65536
)

const (

	// Types of messages
	// The server and the client send messages over the network to
	// communicate with each other. The messages can be of any type which
	// includes audio messages, video messages, command messages, shared
	// object messages, data messages, and user control messages.
	// 3.1. Command message
	// Command messages carry the AMF-encoded commands between the client
	// and the server. These messages have been assigned message type value
	// of 20 for AMF0 encoding and message type value of 17 for AMF3
	// encoding. These messages are sent to perform some operations like
	// connect, createStream, publish, play, pause on the peer. Command
	// messages like onstatus, result etc. are used to inform the sender
	// about the status of the requested commands. A command message
	// consists of command name, transaction ID, and command object that
	// contains related parameters. A client or a server can request Remote
	// Procedure Calls (RPC) over streams that are communicated using the
	// command messages to the peer.

	// RtmpMsgAmf3CommandMessage .
	RtmpMsgAmf3CommandMessage = 17 // 0x11

	// RtmpMsgAmf0CommandMessage .
	RtmpMsgAmf0CommandMessage = 20 // 0x14

	// Data message
	// The client or the server sends this message to send Metadata or any
	// user data to the peer. Metadata includes details about the
	// data(audio, video etc.) like creation time, duration, theme and so
	// on. These messages have been assigned message type value of 18 for
	// AMF0 and message type value of 15 for AMF3.

	// RtmpMsgAmf0DataMessage .
	RtmpMsgAmf0DataMessage = 18 // 0x12

	// RtmpMsgAmf3DataMessage .
	RtmpMsgAmf3DataMessage = 15 // 0x0F

	// Shared object message
	// A shared object is a Flash object (a collection of name value pairs)
	// that are in synchronization across multiple clients, instances, and
	// so on. The message types kMsgContainer=19 for AMF0 and
	// kMsgContainerEx=16 for AMF3 are reserved for shared object events.
	// Each message can contain multiple events.

	// RtmpMsgAmf3SharedObject .
	RtmpMsgAmf3SharedObject = 16 // 0x10

	// RtmpMsgAmf0SharedObject .
	RtmpMsgAmf0SharedObject = 19 // 0x13

	// RtmpMsgAudioMessage  Audio message
	// The client or the server sends this message to send audio data to the
	// peer. The message type value of 8 is reserved for audio messages.
	RtmpMsgAudioMessage = 8 // 0x08

	// RtmpMsgVideoMessage Video message
	// The client or the server sends this message to send video data to the
	// peer. The message type value of 9 is reserved for video messages.
	// These messages are large and can delay the sending of other type of
	// messages. To avoid such a situation, the video message is assigned
	// the lowest priority.
	RtmpMsgVideoMessage = 9 // 0x09

	// RtmpMsgAggregateMessage Aggregate message
	// An aggregate message is a single message that contains a list of submessages.
	// The message type value of 22 is reserved for aggregate
	// messages.
	RtmpMsgAggregateMessage = 22 // 0x16

)

const (

	// RTMP reserves message type IDs 1-7 for protocol control messages.
	// These messages contain information needed by the RTM Chunk Stream
	// protocol or RTMP itself. Protocol messages with IDs 1 & 2 are
	// reserved for usage with RTM Chunk Stream protocol. Protocol messages
	// with IDs 3-6 are reserved for usage of RTMP. Protocol message with ID
	// 7 is used between edge server and origin server.

	// RtmpMsgSetChunkSize .
	RtmpMsgSetChunkSize = 0x01
	// RtmpMsgAbortMessage .
	RtmpMsgAbortMessage = 0x02
	// RtmpMsgAcknowledgement .
	RtmpMsgAcknowledgement = 0x03
	// RtmpMsgUserControlMessage .
	RtmpMsgUserControlMessage = 0x04
	// RtmpMsgWindowAcknowledgementSize .
	RtmpMsgWindowAcknowledgementSize = 0x05
	// RtmpMsgSetPeerBandwidth .
	RtmpMsgSetPeerBandwidth = 0x06
	// RtmpMsgEdgeAndOriginServerCommand .
	RtmpMsgEdgeAndOriginServerCommand = 0x07
)

const (
	// RtmpAmf0Number .
	RtmpAmf0Number = 0x00
	// RtmpAmf0Boolean .
	RtmpAmf0Boolean = 0x01
	// RtmpAmf0String .
	RtmpAmf0String = 0x02
	// RtmpAmf0Object .
	RtmpAmf0Object = 0x03
	// RtmpAmf0MovieClip reserved, not supported
	RtmpAmf0MovieClip = 0x04
	// RtmpAMF0Null .
	RtmpAMF0Null = 0x05
	// RtmpAmf0Undefined .
	RtmpAmf0Undefined = 0x06
	// RtmpAmf0Reference .
	RtmpAmf0Reference = 0x07
	// RtmpAmf0EcmaArray .
	RtmpAmf0EcmaArray = 0x08
	// RtmpAmf0ObjectEnd .
	RtmpAmf0ObjectEnd = 0x09
	// RtmpAmf0StrictArray .
	RtmpAmf0StrictArray = 0x0A
	// RtmpAmf0Date .
	RtmpAmf0Date = 0x0B
	// RtmpAmf0LongString .
	RtmpAmf0LongString = 0x0C
	// RtmpAmf0UnSupported .
	RtmpAmf0UnSupported = 0x0D
	// RtmpAmf0RecordSet reserved, not supported
	RtmpAmf0RecordSet = 0x0E
	// RtmpAmf0XmlDocument .
	RtmpAmf0XmlDocument = 0x0F
	// RtmpAmf0TypedObject .
	RtmpAmf0TypedObject = 0x10
	// RtmpAmf0AVMplusObject AVM+ object is the AMF3 object.
	RtmpAmf0AVMplusObject = 0x11
	// RtmpAmf0OriginStrictArray origin array whos data takes the same form as LengthValueBytes
	RtmpAmf0OriginStrictArray = 0x20
	// RtmpAmf0Invalid User defined
	RtmpAmf0Invalid = 0x3F
)

const (
	// RtmpAmf0CommandConnect .
	RtmpAmf0CommandConnect = "connect"
	// RtmpAmf0CommandCreateStream .
	RtmpAmf0CommandCreateStream = "createStream"
	// RtmpAmf0CommandCloseStream .
	RtmpAmf0CommandCloseStream = "closeStream"
	// RtmpAmf0CommandPlay .
	RtmpAmf0CommandPlay = "play"
	// RtmpAmf0CommandPause .
	RtmpAmf0CommandPause = "pause"
	// RtmpAmf0CommandOnBwDone .
	RtmpAmf0CommandOnBwDone = "onBWDone"
	// RtmpAmf0CommandOnStatus .
	RtmpAmf0CommandOnStatus = "onStatus"
	// RtmpAmf0CommandResult .
	RtmpAmf0CommandResult = "_result"
	// RtmpAmf0CommandError .
	RtmpAmf0CommandError = "_error"
	// RtmpAmf0CommandReleaseStream .
	RtmpAmf0CommandReleaseStream = "releaseStream"
	// RtmpAmf0CommandFcPublish .
	RtmpAmf0CommandFcPublish = "FCPublish"
	// RtmpAmf0CommandUnpublish .
	RtmpAmf0CommandUnpublish = "FCUnpublish"
	// RtmpAmf0CommandPublish .
	RtmpAmf0CommandPublish = "publish"
	// RtmpAmf0CommandGetStreamLength .
	RtmpAmf0CommandGetStreamLength = "getStreamLength"
	// RtmpAmf0CommandKeeplive .
	RtmpAmf0CommandKeeplive = "JMS.KeepAlive"
	// RtmpAmf0CommandEnableVideo .
	RtmpAmf0CommandEnableVideo = "JMS.EnableVideo"
	// RtmpAmf0CommandInsertKeyFrame .
	RtmpAmf0CommandInsertKeyFrame = "JMS.InsertKeyframe"
	// RtmpAmf0DataSampleAccess .
	RtmpAmf0DataSampleAccess = "|RtmpSampleAccess"
	// RtmpAmf0DataSetDataFrame .
	RtmpAmf0DataSetDataFrame = "@setDataFrame"
	// RtmpAmf0DataOnMetaData .
	RtmpAmf0DataOnMetaData = "onMetaData"
	// RtmpAmf0DataOnCustomData .
	RtmpAmf0DataOnCustomData = "onCustomData"
)

const (
	// StatusLevel .
	StatusLevel = "level"
	// StatusCode .
	StatusCode = "code"
	// StatusDescription .
	StatusDescription = "description"
	// StatusDetails .
	StatusDetails = "details"
	// StatusClientID .
	StatusClientID = "clientid"
	// StatusLevelStatus .
	StatusLevelStatus = "status"
	// StatusLevelError status error
	StatusLevelError = "error"

	// StatusCodeConnectSuccess .
	StatusCodeConnectSuccess = "NetConnection.Connect.Success"
	// StatusCodeConnectRejected .
	StatusCodeConnectRejected = "NetConnection.Connect.Rejected"
	// StatusCodeStreamReset .
	StatusCodeStreamReset = "NetStream.Play.Reset"
	// StatusCodeStreamStart .
	StatusCodeStreamStart = "NetStream.Play.Start"
	// StatusCodeStreamPause .
	StatusCodeStreamPause = "NetStream.Pause.Notify"
	// StatusCodeStreamUnpause .
	StatusCodeStreamUnpause = "NetStream.Unpause.Notify"
	// StatusCodePublishStart .
	StatusCodePublishStart = "NetStream.Publish.Start"
	// StatusCodeDataStart .
	StatusCodeDataStart = "NetStream.Data.Start"
	// StatusCodeUnpublishSuccess .
	StatusCodeUnpublishSuccess = "NetStream.Unpublish.Success"

	// FMLE

	// RtmpAmf0CommandOnFcPublish .
	RtmpAmf0CommandOnFcPublish = "onFCPublish"
	// RtmpAmf0CommandOnFcUnpublish .
	RtmpAmf0CommandOnFcUnpublish = "onFCUnpublish"
)

const (

	// band width check method name, which will be invoked by client.
	// band width check mothods use SrsBandwidthPacket as its internal packet type,
	// so ensure you set command name when you use it.
	// server play control

	// RtmpBwCheckStartPlay .
	RtmpBwCheckStartPlay = "onSrsBandCheckStartPlayBytes"
	// RtmpBwCheckStartingPlay .
	RtmpBwCheckStartingPlay = "onSrsBandCheckStartingPlayBytes"
	// RtmpBwCheckStopPlay .
	RtmpBwCheckStopPlay = "onSrsBandCheckStopPlayBytes"
	// RtmpBwCheckStoppedPlay .
	RtmpBwCheckStoppedPlay = "onSrsBandCheckStoppedPlayBytes"

	// server publish control

	// RtmpBwCheckStartPublish .
	RtmpBwCheckStartPublish = "onSrsBandCheckStartPublishBytes"
	// RtmpBwCheckStartingPublish .
	RtmpBwCheckStartingPublish = "onSrsBandCheckStartingPublishBytes"
	// RtmpBwCheckStopPublish .
	RtmpBwCheckStopPublish = "onSrsBandCheckStopPublishBytes"
	// RtmpBwCheckStopppedPublish flash never send out this packet, for its queue is full.
	RtmpBwCheckStopppedPublish = "onSrsBandCheckStoppedPublishBytes"

	// EOF control.

	// RtmpBwCheckFinished the report packet when check finished.
	RtmpBwCheckFinished = "onSrsBandCheckFinished"
	// RtmpBwCheckFinal flash never send out this packet, for its queue is full.
	RtmpBwCheckFinal = "finalClientPacket"

	// data packets

	// RtmpBwCheckPlaying .
	RtmpBwCheckPlaying = "onSrsBandCheckPlaying"
	// RtmpBwCheckPublishing .
	RtmpBwCheckPublishing = "onSrsBandCheckPublishing"
)

const (
	// SrcPCUCStreamBegin The server sends this event to notify the client
	// that a stream has become functional and can be
	// used for communication. By default, this event
	// is sent on ID 0 after the application connect
	// command is successfully received from the
	// client. The event data is 4-byte and represents
	// the stream ID of the stream that became
	// functional.
	SrcPCUCStreamBegin = 0x00

	// SrcPCUCStreamEOF The server sends this event to notify the client
	// that the playback of data is over as requested
	// on this stream. No more data is sent without
	// issuing additional commands. The client discards
	// the messages received for the stream. The
	// 4 bytes of event data represent the ID of the
	// stream on which playback has ended.
	SrcPCUCStreamEOF = 0x01

	// SrcPCUCStreamDry The server sends this event to notify the client
	// that there is no more data on the stream. If the
	// server does not detect any message for a time
	// period, it can notify the subscribed clients
	// that the stream is dry. The 4 bytes of event
	// data represent the stream ID of the dry stream.
	SrcPCUCStreamDry = 0x02

	// SrcPCUCSetBufferLength The client sends this event to inform the server
	// of the buffer size (in milliseconds) that is
	// used to buffer any data coming over a stream.
	// This event is sent before the server starts
	// processing the stream. The first 4 bytes of the
	// event data represent the stream ID and the next
	// 4 bytes represent the buffer length, in
	// milliseconds.
	SrcPCUCSetBufferLength = 0x03 // 8bytes event-data

	// SrcPCUCStreamIsRecorded The server sends this event to notify the client
	// that the stream is a recorded stream. The
	// 4 bytes event data represent the stream ID of
	// the recorded stream.
	SrcPCUCStreamIsRecorded = 0x04

	// SrcPCUCPingRequest The server sends this event to test whether the
	// client is reachable. Event data is a 4-byte
	// timestamp, representing the local server time
	// when the server dispatched the command. The
	// client responds with kMsgPingResponse on
	// receiving kMsgPingRequest.
	SrcPCUCPingRequest = 0x06

	// SrcPCUCPingResponse The client sends this event to the server in
	// response to the ping request. The event data is
	// a 4-byte timestamp, which was received with the
	// kMsgPingRequest request.
	SrcPCUCPingResponse = 0x07
)

const (
	// RtmpSigClientID signature for packets to client
	RtmpSigClientID = "ASAICiss"
	// RtmpSigFmsVersion signature for packets to client
	RtmpSigFmsVersion = "1.0.0.0"
	// RtmpSigAmf0Ver objectEncoding default value.
	RtmpSigAmf0Ver = 0.0
)

// flv format
// AACPacketType IF SoundFormat == 10 UI8
// The following values are defined:
//     0 = AAC sequence header
//     1 = AAC raw
const (
	// RtmpCodecAudioTypeReserved set to the max value to reserved, for array map.
	RtmpCodecAudioTypeReserved = 2

	// RtmpCodecAudioTypeSequenceHeader audio type sequence header
	RtmpCodecAudioTypeSequenceHeader = 0
	// RtmpCodecAudioTypeRawData audio raw data
	RtmpCodecAudioTypeRawData = 1
)

// E.4.3.1 VIDEODATA
// Frame Type UB [4]
// Type of video frame. The following values are defined:
//     1 = key frame (for AVC, a seekable frame)
//     2 = inter frame (for AVC, a non-seekable frame)
//     3 = disposable inter frame (H.263 only)
//     4 = generated key frame (reserved for server use only)
//     5 = video info/command frame
const (
	// RtmpCodecVideoAVCFrameReserved set to the max value to reserved, for array map.
	RtmpCodecVideoAVCFrameReserved = 0
	// RtmpCodecVideoAVCFrameReserved1 .
	RtmpCodecVideoAVCFrameReserved1 = 6

	// RtmpCodecVideoAVCFrameKeyFrame video h264 key frame
	RtmpCodecVideoAVCFrameKeyFrame = 1
	// RtmpCodecVideoAVCFrameInterFrame video h264 inter frame
	RtmpCodecVideoAVCFrameInterFrame = 2
	// RtmpCodecVideoAVCFrameDisposableInterFrame .
	RtmpCodecVideoAVCFrameDisposableInterFrame = 3
	// RtmpCodecVideoAVCFrameGeneratedKeyFrame .
	RtmpCodecVideoAVCFrameGeneratedKeyFrame = 4
	// RtmpCodecVideoAVCFrameVideoInfoFrame .
	RtmpCodecVideoAVCFrameVideoInfoFrame = 5
)

// AVCPacketType IF CodecID == 7 UI8
// The following values are defined:
//     0 = AVC sequence header
//     1 = AVC NALU
//     2 = AVC end of sequence (lower level NALU sequence ender is
//         not required or supported)
const (
	// set to the max value to reserved, for array map.
	RtmpCodecVideoAVCTypeReserved = 3

	// RtmpCodecVideoAVCTypeSequenceHeader .
	RtmpCodecVideoAVCTypeSequenceHeader = 0
	// RtmpCodecVideoAVCTypeNALU .
	RtmpCodecVideoAVCTypeNALU = 1
	// RtmpCodecVideoAVCTypeSequenceHeaderEOF .
	RtmpCodecVideoAVCTypeSequenceHeaderEOF = 2
)

// E.4.3.1 VIDEODATA
// CodecID UB [4]
// Codec Identifier. The following values are defined:
//     2 = Sorenson H.263
//     3 = Screen video
//     4 = On2 VP6
//     5 = On2 VP6 with alpha channel
//     6 = Screen video version 2
//     7 = AVC
//     13 = HEVC
const (
	// RtmpCodecVideoReserved set to the max value to reserved, for array map.
	RtmpCodecVideoReserved = 0
	// RtmpCodecVideoReserved1 .
	RtmpCodecVideoReserved1 = 1
	// RtmpCodecVideoReserved2 .
	RtmpCodecVideoReserved2 = 8

	// RtmpCodecVideoSorensonH263 .
	RtmpCodecVideoSorensonH263 = 2
	// RtmpCodecVideoScreenVideo .
	RtmpCodecVideoScreenVideo = 3
	// RtmpCodecVideoOn2VP6 .
	RtmpCodecVideoOn2VP6 = 4
	// RtmpCodecVideoOn2VP6WithAlphaChannel .
	RtmpCodecVideoOn2VP6WithAlphaChannel = 5
	// RtmpCodecVideoScreenVideoVersion2 .
	RtmpCodecVideoScreenVideoVersion2 = 6
	// RtmpCodecVideoAVC .
	RtmpCodecVideoAVC = 7
	// RtmpCodecVideoHEVC h265
	RtmpCodecVideoHEVC = 13
)

// SoundFormat UB [4]
// Format of SoundData. The following values are defined:
//     0 = Linear PCM, platform endian
//     1 = ADPCM
//     2 = MP3
//     3 = Linear PCM, little endian
//     4 = Nellymoser 16 kHz mono
//     5 = Nellymoser 8 kHz mono
//     6 = Nellymoser
//     7 = G.711 A-law logarithmic PCM
//     8 = G.711 mu-law logarithmic PCM
//     9 = reserved
//     10 = AAC
//     11 = Speex
//     14 = MP3 8 kHz
//     15 = Device-specific sound
// Formats 7, 8, 14, and 15 are reserved.
// AAC is supported in Flash Player 9,0,115,0 and higher.
// Speex is supported in Flash Player 10 and higher.
const (
	// RtmpCodecAudioReserved1 set to the max value to reserved, for array map.
	RtmpCodecAudioReserved1 = 16

	// RtmpCodecAudioLinearPCMPlatformEndian .
	RtmpCodecAudioLinearPCMPlatformEndian = 0
	// RtmpCodecAudioADPCM .
	RtmpCodecAudioADPCM = 1
	// RtmpCodecAudioMP3 .
	RtmpCodecAudioMP3 = 2
	// RtmpCodecAudioLinearPCMLittleEndian .
	RtmpCodecAudioLinearPCMLittleEndian = 3
	// RtmpCodecAudioNellymoser16kHzMono .
	RtmpCodecAudioNellymoser16kHzMono = 4
	// RtmpCodecAudioNellymoser8kHzMono .
	RtmpCodecAudioNellymoser8kHzMono = 5
	// RtmpCodecAudioNellymoser .
	RtmpCodecAudioNellymoser = 6
	// RtmpCodecAudioReservedG711AlawLogarithmicPCM .
	RtmpCodecAudioReservedG711AlawLogarithmicPCM = 7
	// RtmpCodecAudioReservedG711MuLawLogarithmicPCM .
	RtmpCodecAudioReservedG711MuLawLogarithmicPCM = 8
	// RtmpCodecAudioReserved .
	RtmpCodecAudioReserved = 9
	// RtmpCodecAudioAAC .
	RtmpCodecAudioAAC = 10
	// RtmpCodecAudioSpeex .
	RtmpCodecAudioSpeex = 11
	// RtmpCodecAudioReservedMP3Of8kHz .
	RtmpCodecAudioReservedMP3Of8kHz = 14
	// RtmpCodecAudioReservedDeviceSpecificSound .
	RtmpCodecAudioReservedDeviceSpecificSound = 15
)

// TokenStr token for auth
const TokenStr = "?token="
