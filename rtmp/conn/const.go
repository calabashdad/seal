package conn

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
