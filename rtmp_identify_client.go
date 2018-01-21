package main

func (rtmp *RtmpSession) IdendifyClient() (err error) {

	for {
		var chunk *ChunkStream

		err, chunk = rtmp.RecvMsg()
		if err != nil {
			return
		}
		err = rtmp.DecodeMsg(chunk)
		if err != nil {
			return
		}
	}

	return
}
