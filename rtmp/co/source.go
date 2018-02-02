package co

import (
	"log"
	"sync"
)

type SourceInfoHub struct {
	//key: streamName
	hub  map[string]*SourceInfoS
	lock sync.RWMutex
}

var sourcesHub = &SourceInfoHub{
	hub: make(map[string]*SourceInfoS),
}

//data source info.
type SourceInfoS struct {
	sampleRate float64 // the sample rate of audio in metadata
	frameRate  float64 // the video frame rate in metadata
	atc        bool    // atc whether atc(use absolute time and donot adjust time),
	// directly use msg time and donot adjust if atc is true,
	// otherwise, adjust msg time to start from 0 to make flash happy.

}

func (sh *SourceInfoHub) findSourceToPublish(k string) (s *SourceInfoS) {

	sh.lock.Lock()
	defer sh.lock.Unlock()

	res := sh.hub[k]
	if nil != res {
		log.Println("stream ", k, " can not publish, because has already publishing....")
		return nil
	}

	//can publish. new a source
	sh.hub[k] = &SourceInfoS{}

	return sh.hub[k]
}

func (sh *SourceInfoHub) findSourceToPlay(k string) (s *SourceInfoS) {
	sh.lock.Lock()
	defer sh.lock.Unlock()

	return nil
}

func (sh *SourceInfoHub) deleteSource(streamName string) {
	sh.lock.Lock()
	defer sh.lock.Unlock()

	delete(sh.hub, streamName)
}
