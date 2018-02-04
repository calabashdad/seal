package co

import (
	"fmt"
	"log"
	"seal/rtmp/pt"
	"time"
)

//Consumer is the consumer of source
type Consumer struct {
	queueSizeMills uint32
	avStartTime    int64
	avEndTime      int64
	msgs           chan *pt.Message
	jitter         struct {
		lastPktTime        int64
		lastPktCorrectTime int64
	}
	paused bool
}

// atc whether atc, donot use jitter correct if true
// tba timebase of audio. used to calc the audio time delta if time-jitter detected.
// tbv timebase of video. used to calc the video time delta if time-jitter detected.
func (c *Consumer) enquene(msg *pt.Message, atc bool, tba float64, tbv float64, timeJitter uint32) {

	if nil == msg {
		return
	}

	if !atc {
		c.timeJittrCorrect(msg, tba, tbv, timeJitter)
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

func (c *Consumer) dump() (err error, msg *pt.Message) {

	select {
	case <-time.After(time.Duration(5) * time.Second):
		err = fmt.Errorf("wait 3 seconds, the source is dry")
		return
	case msg = <-c.msgs:
		return
	}

	return
}

func (c *Consumer) timeJittrCorrect(msg *pt.Message, tba float64, tbv float64, timeJitter uint32) {

	if nil == msg {
		return
	}

	if RtmpTimeJitterFull != timeJitter {
		// all jitter correct features is disabled, ignore.
		if RtmpTimeJitterOff == timeJitter {
			return
		}

		// start at zero, but donot ensure monotonically increasing.
		if RtmpTimeJitterZero == timeJitter {
			// for the first time, last_pkt_correct_time is zero.
			// while when timestamp overflow, the timestamp become smaller,
			// reset the last_pkt_correct_time.
			if c.jitter.lastPktCorrectTime <= 0 || c.jitter.lastPktCorrectTime > int64(msg.Header.Timestamp) {
				c.jitter.lastPktCorrectTime = int64(msg.Header.Timestamp)
			}

			msg.Header.Timestamp -= uint64(c.jitter.lastPktCorrectTime)

			return
		}
	}

	// full jitter algorithm, do jitter correct.

	// set to 0 for metadata.
	if !msg.Header.IsAudio() && !msg.Header.IsVideo() {
		msg.Header.Timestamp = 0
		return
	}

	sampleRate := tba
	frameRate := tbv

	/**
	 * we use a very simple time jitter detect/correct algorithm:
	 * 1. delta: ensure the delta is positive and valid,
	 *     we set the delta to DEFAULT_FRAME_TIME_MS,
	 *     if the delta of time is nagative or greater than CONST_MAX_JITTER_MS.
	 * 2. last_pkt_time: specifies the original packet time,
	 *     is used to detect next jitter.
	 * 3. last_pkt_correct_time: simply add the positive delta,
	 *     and enforce the time monotonically.
	 */
	timeLocal := msg.Header.Timestamp
	delta := int64(timeLocal) - c.jitter.lastPktTime

	// if jitter detected, reset the delta.
	if delta < 0 || delta > CONST_MAX_JITTER_MS {
		// calc the right diff by audio sample rate
		if msg.Header.IsAudio() && sampleRate > 0 {
			delta = (int64)(float64(delta) * 1000.0 / sampleRate)
		} else if msg.Header.IsVideo() && frameRate > 0 {
			delta = (int64)(float64(delta) * 1.0 / frameRate)
		} else {
			delta = DEFAULT_FRAME_TIME_MS
		}
	}

	// sometimes, the time is absolute time, so correct it again.
	if delta < 0 || delta > CONST_MAX_JITTER_MS {
		delta = DEFAULT_FRAME_TIME_MS
	}

	if c.jitter.lastPktCorrectTime+delta > 0 {
		c.jitter.lastPktCorrectTime = c.jitter.lastPktCorrectTime + delta
	} else {
		c.jitter.lastPktCorrectTime = 0
	}

	msg.Header.Timestamp = uint64(c.jitter.lastPktCorrectTime)
	c.jitter.lastPktCorrectTime = int64(timeLocal)

}

func (c *Consumer) onPlayPause(isPause bool) (err error) {
	c.paused = true
	log.Println("consumer changed pause status to ", isPause)
	return
}
