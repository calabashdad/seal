package co

import (
	"seal/rtmp/pt"
)

type hlsStream struct {
}

func newHlsStream() *hlsStream {
	return &hlsStream{}
}

func (hls *hlsStream) onMeta(pkt *pt.OnMetaDataPacket) (err error) {
	return
}

func (hls *hlsStream) onAudio(msg *pt.Message) (err error) {
	return
}

func (hls *hlsStream) onVideo(msg *pt.Message) (err error) {
	return
}
