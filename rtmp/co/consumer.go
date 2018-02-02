package co

import (
	"seal/rtmp/pt"
	"time"
)

//Consumer is the consumer of source
type Consumer struct {
	queueSizeMills uint32
	avStartTime    int64
	avEndTime      int64
	msgs           chan *pt.Message
}

// atc whether atc, donot use jitter correct if true
// tba timebase of audio. used to calc the audio time delta if time-jitter detected.
// tbv timebase of video. used to calc the video time delta if time-jitter detected.
func (c *Consumer) enquene(msg *pt.Message, atc bool, tba float64, tbv float64, timeJitter uint32) {

	if !atc {
		//todo. time jitter.
	}

	if msg.Header.IsVideo() || msg.Header.IsAudio() {
		if -1 == c.avStartTime {
			c.avStartTime = int64(msg.Header.Timestamp)
		}

		c.avEndTime = int64(msg.Header.Timestamp)
	}

	//push into chan
	select {
	//add timeout in case block there
	case <-time.After(2 * time.Millisecond):
	case c.msgs <- msg:
	}

	//shrink
	for {
		if uint32(c.avEndTime-c.avStartTime) < c.queueSizeMills {
			break
		}

		//this may be a bug, when msg is video key frame, do not pop it.
		v := <-c.msgs
		c.avStartTime = int64(v.Header.Timestamp)

	}
}
