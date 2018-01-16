package main

func (rtmp *RtmpSession) Connect() (err error) {

	err = rtmp.ExpectMsg()
	if err != nil {
		return
	}

	return
}
