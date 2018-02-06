package co

import (
	"log"
	"seal/rtmp/flv"
	"seal/rtmp/pt"
	"sync"
)

type messagesQuene struct {
	msgs []*pt.Message
	mu   sync.RWMutex
}

func (q *messagesQuene) queue(msg *pt.Message) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.msgs = append(q.msgs, msg)
}

func (c *Consumer) shink() {

	c.msgQuene.mu.Lock()
	defer c.msgQuene.mu.Unlock()

	var iFrameIndex int = -1
	for i := 1; i < len(c.msgQuene.msgs); i++ {

		if uint32(c.avEndTime-c.avStartTime) < c.queueSizeMills {
			break
		}

		if c.msgQuene.msgs[i].Header.IsVideo() {
			if flv.VideoH264IsSpspps(c.msgQuene.msgs[i].Payload.Payload) {
				//the max iframe index to remove
				iFrameIndex = i

				// set the start time, we will remove until this frame.
				c.avStartTime = int64(c.msgQuene.msgs[i].Header.Timestamp)

				break
			}
		}

		c.avStartTime = int64(c.msgQuene.msgs[i].Header.Timestamp)
	}

	// no iframe, for audio, clear the queue.
	// it is ok to clear for audio, for the shrink tell us the queue is full.
	// for video, we clear util the I-Frame, for the decoding must start from I-frame,
	// for audio, it's ok to clear any data, also we can clear the whole queue.
	if iFrameIndex < 0 {
		//clear
		c.msgQuene.msgs = c.msgQuene.msgs[:0]
	} else {
		c.msgQuene.msgs = append(c.msgQuene.msgs[:iFrameIndex], c.msgQuene.msgs[iFrameIndex:]...)
	}
}

//Consumer is the consumer of source
type Consumer struct {
	queueSizeMills uint32
	avStartTime    int64
	avEndTime      int64
	msgQuene       messagesQuene
	jitter         struct {
		lastPktTime        int64
		lastPktCorrectTime int64
	}
	paused   bool
	Duration float64
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

	c.msgQuene.queue(msg)

	log.Println("equene an msg, msg type=", msg.Header.MessageType,
		",stream id =", msg.Header.StreamId,
		",timestamp=", msg.Header.Timestamp,
		",payload=", len(msg.Payload.Payload))

	c.shink()
}

func (c *Consumer) dump() (err error, msg []*pt.Message) {

	c.msgQuene.mu.Lock()
	defer c.msgQuene.mu.Unlock()

	if c.paused {
		return
	}

	if 0 == len(c.msgQuene.msgs) {
		return
	}

	var countOnceDump int = -1

	if len(c.msgQuene.msgs) > 30 {
		countOnceDump = 30
	} else {
		countOnceDump = len(c.msgQuene.msgs)
	}

	for i := 0; i < countOnceDump; i++ {
		msg = append(msg, c.msgQuene.msgs[i])
	}

	c.msgQuene.msgs = append(c.msgQuene.msgs[:countOnceDump], c.msgQuene.msgs[countOnceDump:]...)

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
	c.paused = isPause
	log.Println("consumer changed pause status to ", isPause)
	return
}
