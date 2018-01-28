package conn

const (
	/**
	* max rtmp header size:
	*     1bytes basic header,
	*     11bytes message header,
	*     4bytes timestamp header,
	* that is, 1+11+4=16bytes.
	 */
	RTMP_MAX_FMT0_HEADER_SIZE = 16
)
