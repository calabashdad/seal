package co

import (
	"log"
	"seal/conf"
	"seal/rtmp/pt"
	"sync"
)

type SourceHub struct {
	//key: streamName
	hub  map[string]*Source
	lock sync.RWMutex
}

var sourcesHub = &SourceHub{
	hub: make(map[string]*Source),
}

//data source info.
type Source struct {
	// the sample rate of audio in metadata
	sampleRate float64
	// the video frame rate in metadata
	frameRate float64
	// atc whether atc(use absolute time and donot adjust time),
	// directly use msg time and donot adjust if atc is true,
	// otherwise, adjust msg time to start from 0 to make flash happy.
	atc bool
	//time jitter algrithem
	timeJitter uint32
	//cached meta data
	cacheMetaData *pt.Message
	//cached video sequence header
	cacheVideoSequenceHeader *pt.Message
	//cached aideo sequence header
	cacheAudioSequenceHeader *pt.Message

	//consumers
	consumers map[*Consumer]*Consumer
	//lock for consumers.
	consumerLock sync.RWMutex

	//gop cache
	gopCache *GopCache
}

func (s *Source) CreateConsumer(c *Consumer) {
	if nil == c {
		log.Println("when registe consumer, nil == consumer")
		return
	}

	s.consumerLock.Lock()
	defer s.consumerLock.Unlock()

	s.consumers[c] = c

}

func (s *Source) DestroyConsumer(c *Consumer) {
	if nil == c {
		log.Println("when destroy consumer, nil == consummer")
		return
	}

	s.consumerLock.Lock()
	defer s.consumerLock.Unlock()

	delete(s.consumers, c)

}

func (s *Source) copyToAllConsumers(msg *pt.Message) {
	s.consumerLock.Lock()
	defer s.consumerLock.Unlock()

	for _, v := range s.consumers {
		v.enquene(msg, s.atc, s.sampleRate, s.frameRate, s.timeJitter)
	}
}

func (s *SourceHub) findSourceToPublish(k string) *Source {

	s.lock.Lock()
	defer s.lock.Unlock()

	res := s.hub[k]
	if nil != res {
		log.Println("stream ", k, " can not publish, because has already publishing....")
		return nil
	}

	//can publish. new a source
	s.hub[k] = &Source{
		timeJitter: conf.GlobalConfInfo.Rtmp.TimeJitter,
		gopCache:   &GopCache{},
	}

	return s.hub[k]
}

func (s *SourceHub) findSourceToPlay(k string) *Source {
	s.lock.Lock()
	defer s.lock.Unlock()

	res := s.hub[k]
	if nil != res {
		return res
	}

	log.Println("stream ", k, " can not play, because has not published.")

	return nil
}

func (s *SourceHub) deleteSource(streamName string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.hub, streamName)
}
