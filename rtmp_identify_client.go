package main

func (rtmp *RtmpSession) IdendifyClient() (err error) {

	for {
		var chunk *ChunkStream

		err, chunk = rtmp.RecvMsg()
		if err != nil {
			break
		}
		err = rtmp.DecodeMsg(chunk)
		if err != nil {
			break
		}

		err = rtmp.HandleMsg(chunk)
		if err != nil {
			break
		}
	}

	if err != nil {
		return
	}

	return
}
