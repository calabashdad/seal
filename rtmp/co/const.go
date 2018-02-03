package co

const (
	//RtmpMaxFmt0HeaderSize is the max rtmp header size:
	//   1bytes basic header,
	//   11bytes message header,
	//   4bytes timestamp header,
	//   that is, 1+11+4=16bytes.
	RtmpMaxFmt0HeaderSize = 16
)

//RtmpRole define
const (
	RtmpRoleUnknown = 0

	RtmpRoleFMLEPublisher  = 1
	RtmpRoleFlashPublisher = 2
	RtmpRolePlayer         = 3
)

const (
	RtmpTimeJitterFull = 0x01
	RtmpTimeJitterZero = 0x02
	RtmpTimeJitterOff  = 0x03
)

// for 26ms per audio packet,
// 115 packets is 3s.
const (
	PURE_AUDIO_GUESS_COUNT = 115
)

const (
	CONST_MAX_JITTER_MS   = 500
	DEFAULT_FRAME_TIME_MS = 40
)
