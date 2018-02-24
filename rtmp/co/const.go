package co

const (
	// RtmpMaxFmt0HeaderSize is the max rtmp header size:
	//   1bytes basic header,
	//   11bytes message header,
	//   4bytes timestamp header,
	//   that is, 1+11+4=16bytes.
	RtmpMaxFmt0HeaderSize = 16
)

//RtmpRole define
const (
	// RtmpRoleUnknown role undefined
	RtmpRoleUnknown = 0
	// RtmpRoleFMLEPublisher role fmle publisher
	RtmpRoleFMLEPublisher = 1
	// RtmpRoleFlashPublisher role flash publisher
	RtmpRoleFlashPublisher = 2
	// RtmpRolePlayer role player
	RtmpRolePlayer = 3
)

const (
	// RtmpTimeJitterFull time jitter full mode, to ensure stream start at zero, and ensure stream monotonically increasing.
	RtmpTimeJitterFull = 0x01
	// RtmpTimeJitterZero zero mode, only ensure sttream start at zero, ignore timestamp jitter.
	RtmpTimeJitterZero = 0x02
	// RtmpTimeJitterOff off mode, disable the time jitter algorithm, like atc.
	RtmpTimeJitterOff = 0x03
)

const (
	// PureAudioGuessCount for 26ms per audio packet,
	// 115 packets is 3s.
	PureAudioGuessCount = 115
)

const (
	// MaxJitterMs max time delta, which is the between localtime and last packet time
	MaxJitterMs = 500
	// DefaultFrameTimeMs default time delta, which is the between localtime and last packet time
	DefaultFrameTimeMs = 40
)
