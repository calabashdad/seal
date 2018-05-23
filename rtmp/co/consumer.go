package co

import (
	"log"
	"seal/conf"
	"seal/rtmp/pt"
	"time"
)

//Consumer is the consumer of source
type Consumer struct {
	stream         string
	queueSizeMills uint32
	avStartTime    int64
	avEndTime      int64
	msgQuene       chan *pt.Message
	jitter         *pt.TimeJitter
	paused         bool
	duration       float64
}

func NewConsumer(key string) *Consumer {
	return &Consumer{
		stream:         key,
		queueSizeMills: conf.GlobalConfInfo.Rtmp.ConsumerQueueSize * 1000,
		avStartTime:    -1,
		avEndTime:      -1,
		msgQuene:       make(chan *pt.Message, 4096),
		jitter:         &pt.TimeJitter{},
	}
}

func (c *Consumer) Clean() {
	close(c.msgQuene)
}

// Atc whether Atc, donot use jitter correct if true
// tba timebase of audio. used to calc the audio time delta if time-jitter detected.
// tbv timebase of video. used to calc the video time delta if time-jitter detected.
func (c *Consumer) Enquene(msg *pt.Message, atc bool, tba float64, tbv float64, timeJitter uint32) {

	if nil == msg {
		return
	}

	if !atc {
		//c.jitter.Correct(msg, tba, tbv, timeJitter)
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
		log.Println("enquene to channel timeout, channel may be full, key=", c.stream)
		break
	case c.msgQuene <- msg:
		break
	}
}

func (c *Consumer) Dump() (msg *pt.Message) {

	if c.paused {
		log.Println("client paused now")
		return
	}

	select {
	case <-time.After(time.Duration(10) * time.Millisecond):
		// in case block
		return
	case msg = <-c.msgQuene:
		break
	}

	//for {
	//	var msgLocal *pt.Message
	//
	//	select {
	//	case <-time.After(time.Duration(10) * time.Millisecond):
	//		// in case block
	//		return
	//	case msgLocal = <-c.msgQuene:
	//		break
	//	}
	//
	//	c.avStartTime = int64(msgLocal.Header.Timestamp)
	//
	//	if uint32(c.avEndTime-c.avStartTime) > c.queueSizeMills {
	//
	//		// for metadata, key frame, audio sequence header, we do not shrink it.
	//		if msgLocal.Header.IsAmf0Data() ||
	//			flv.VideoH264IsKeyframe(msgLocal.Payload.Payload) ||
	//			flv.AudioIsSequenceHeader(msgLocal.Payload.Payload) {
	//
	//			msg = msgLocal
	//			log.Printf("key=%s dump a frame even it's timestamp is too old. msg type=%d, msg time=%d, avStatrtTime=%d, queue len=%d\n",
	//				c.stream, msg.Header.MessageType, msgLocal.Header.Timestamp, c.avStartTime, c.queueSizeMills)
	//
	//			return
	//		} else {
	//			// msg is too old, drop it directly, we store the latest i frame into cache already
	//			continue
	//		}
	//	} else {
	//		msg = msgLocal
	//		return
	//	}
	//
	//}

	return
}

func (c *Consumer) onPlayPause(isPause bool) (err error) {
	c.paused = isPause
	log.Println("consumer changed pause status to ", isPause)
	return
}
