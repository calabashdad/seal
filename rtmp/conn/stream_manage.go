package conn

import (
	"sync"
)

//StreamsHub key:stream name value: interface{}
type StreamsHub struct {
	Hub  map[string]interface{}
	Lock sync.Mutex
}

var GlobalStreamsHub StreamsHub = StreamsHub{
	Hub: make(map[string]interface{}),
}

//CheckStreamCanPublish check the stream can publish, one stream can publish unique.
func (rc *RtmpConn) CheckStreamCanPublish(streamName string) bool {

	GlobalStreamsHub.Lock.Lock()
	defer GlobalStreamsHub.Lock.Unlock()

	if res := GlobalStreamsHub.Hub[streamName]; nil != res {
		//the stream has published.
		return false
	}

	//stream can publish
	GlobalStreamsHub.Hub[streamName] = struct{}{}

	return true
}

func (rc *RtmpConn) DeletePublishStream(streamName string) {
	GlobalStreamsHub.Lock.Lock()
	defer GlobalStreamsHub.Lock.Unlock()

	delete(GlobalStreamsHub.Hub, streamName)
}
