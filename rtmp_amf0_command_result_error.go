package main

func (rtmp *RtmpSession) handleAMF0CommandResultError(chunk *ChunkStream) (err error) {

	var offset uint32

	var transactionId float64
	err, transactionId = Amf0ReadNumber(chunk.msg.payload, &offset)
	if err != nil {
		return
	}

	_ = transactionId //todo

	return
}
