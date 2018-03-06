package co

import (
	"log"
	"seal/rtmp/pt"
	"time"
)

//Consumer is the consumer of source
type Consumer struct {
	queueSizeMills uint32
	avStartTime    int64
	avEndTime      int64
	msgQuene       chan *pt.Message
	jitter         *pt.TimeJitter
	paused         bool
	duration       float64
}

// atc whether atc, donot use jitter correct if true
// tba timebase of audio. used to calc the audio time delta if time-jitter detected.
// tbv timebase of video. used to calc the video time delta if time-jitter detected.
func (c *Consumer) enquene(msg *pt.Message, atc bool, tba float64, tbv float64, timeJitter uint32) {

	if nil == msg {
		return
	}

	if !atc {
		c.jitter.Correct(msg, tba, tbv, timeJitter)
	}

	if msg.Header.IsVideo() || msg.Header.IsAudio() {
		if -1 == c.avStartTime {
			c.avStartTime = int64(msg.Header.Timestamp)
		}

		c.avEndTime = int64(msg.Header.Timestamp)
	}

	select {
	// incase block, and influence others.
	case <-time.After(time.Duration(3) * time.Millisecond):
		break
	case c.msgQuene <- msg:
		break
	}
}

func (c *Consumer) dump() (msg *pt.Message) {

	if c.paused {
		log.Println("client paused now")
		return
	}

	select {
	case <-time.After(time.Duration(5) * time.Millisecond):
		// in case block
		break
	case msg = <-c.msgQuene:
		break
	}

	return
}

func (c *Consumer) onPlayPause(isPause bool) (err error) {
	c.paused = isPause
	log.Println("consumer changed pause status to ", isPause)
	return
}
