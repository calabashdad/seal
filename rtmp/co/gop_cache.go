package co

import (
	"log"
	"seal/rtmp/flv"
	"seal/rtmp/pt"
	"sync"
)

// GopCache cache gop of video/audio to enable players fast start
type GopCache struct {
	mu sync.RWMutex

	// cachedVideoCount the video frame count, avoid cache for pure audio stream.
	cachedVideoCount uint32
	// when user disabled video when publishing, and gop cache enalbed,
	// we will cache the audio/video for we already got video, but we never
	// know when to clear the gop cache, for there is no video in future,
	// so we must guess whether user disabled the video.
	// when we got some audios after laster video, for instance, 600 audio packets,
	// about 3s(26ms per packet) 115 audio packets, clear gop cache.
	//
	// it is ok for performance, for when we clear the gop cache,
	// gop cache is disabled for pure audio stream.
	audioAfterLastVideoCount uint32

	msgs []*pt.Message
}

func (g *GopCache) cache(msg *pt.Message) {

	if nil == msg {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	if msg.Header.IsVideo() {
		g.cachedVideoCount++
		g.audioAfterLastVideoCount = 0
	}

	if g.pureAudio() {
		return
	}

	if msg.Header.IsAudio() {
		g.audioAfterLastVideoCount++
	}

	if g.audioAfterLastVideoCount > pt.PureAudioGuessCount {
		g.clear()
		log.Printf("gop cahce pure audio more than %d seconds, clear old caches.\n", pt.PureAudioGuessCount)
		return
	}

	// clear gop cache when got key frame
	if msg.Header.IsVideo() && flv.VideoIsH264(msg.Payload.Payload) && flv.VideoH264IsKeyframe(msg.Payload.Payload) {
		g.clear()

		// curent msg is video frame, so we set to 1.
		g.cachedVideoCount = 1

		log.Println("@@@@@gop cache recv a key frame, clear old caches.")

	}

	g.msgs = append(g.msgs, msg)
}

func (g *GopCache) pureAudio() bool {
	return 0 == g.cachedVideoCount
}

func (g *GopCache) clear() {
	g.msgs = nil
	g.cachedVideoCount = 0
	g.audioAfterLastVideoCount = 0
}

func (g *GopCache) Empty() bool {
	return nil == g.msgs
}

func (g *GopCache) StartTime() uint64 {
	if nil != g.msgs && nil != g.msgs[0] {
		return g.msgs[0].Header.Timestamp
	}

	return 0
}

func (g *GopCache) Dump(c *Consumer, atc bool, tba float64, tbv float64, timeJitter uint32) {

	g.mu.Lock()
	defer g.mu.Unlock()

	for _, v := range g.msgs {

		if nil == v {
			continue
		}

		c.Enquene(v, atc, tba, tbv, timeJitter)
	}
}
