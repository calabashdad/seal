package co

import (
	"log"
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
	sampleRate float64 // the sample rate of audio in metadata
	frameRate  float64 // the video frame rate in metadata
	atc        bool    // atc whether atc(use absolute time and donot adjust time),
	// directly use msg time and donot adjust if atc is true,
	// otherwise, adjust msg time to start from 0 to make flash happy.

}

func (sh *SourceHub) findSourceToPublish(k string) (s *Source) {

	sh.lock.Lock()
	defer sh.lock.Unlock()

	res := sh.hub[k]
	if nil != res {
		log.Println("stream ", k, " can not publish, because has already publishing....")
		return nil
	}

	//can publish. new a source
	sh.hub[k] = &Source{}

	return sh.hub[k]
}

func (sh *SourceHub) findSourceToPlay(k string) (s *Source) {
	sh.lock.Lock()
	defer sh.lock.Unlock()

	res := sh.hub[k]
	if nil != res {
		return res
	}

	log.Println("stream ", k, " can not play, because has not published.")

	return nil
}

func (sh *SourceHub) deleteSource(streamName string) {
	sh.lock.Lock()
	defer sh.lock.Unlock()

	delete(sh.hub, streamName)
}
