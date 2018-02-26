package flv

import (
	"seal/rtmp/pt"
)

// AudioIsSequenceHeader judge audio is aac sequence header
func AudioIsSequenceHeader(data []uint8) bool {

	if !audioIsAAC(data) {
		return false
	}

	if len(data) < 2 {
		return false
	}

	aacPacketType := data[1]

	return aacPacketType == pt.RtmpCodecAudioTypeSequenceHeader
}

func audioIsAAC(data []uint8) bool {

	if len(data) < 1 {
		return false
	}

	soundFormat := data[0]
	soundFormat = (soundFormat >> 4) & 0x0f

	return soundFormat == pt.RtmpCodecAudioAAC
}

// VideoIsH264 judge video is h264 sequence header
func VideoIsH264(data []uint8) bool {

	if len(data) < 1 {
		return false
	}

	codecID := data[0]
	codecID &= 0x0f

	return pt.RtmpCodecVideoAVC == codecID
}

// VideoH264IsKeyframe judge video is h264 key frame
func VideoH264IsKeyframe(data []uint8) bool {
	// 2bytes required.
	if len(data) < 2 {
		return false
	}

	frameType := data[0]
	frameType = (frameType >> 4) & 0x0F

	return frameType == pt.RtmpCodecVideoAVCFrameKeyFrame
}

// VideoH264IsSequenceHeaderAndKeyFrame judge video is h264 sequence header and key frame
func VideoH264IsSequenceHeaderAndKeyFrame(data []uint8) bool {
	// sequence header only for h264
	if !VideoIsH264(data) {
		return false
	}

	// 2bytes required.
	if len(data) < 2 {
		return false
	}

	frameType := data[0]
	frameType = (frameType >> 4) & 0x0F

	avcPacketType := data[1]

	return frameType == pt.RtmpCodecVideoAVCFrameKeyFrame && avcPacketType == pt.RtmpCodecVideoAVCTypeSequenceHeader
}

// VideoH264IsSpspps judge video is spspps
func VideoH264IsSpspps(data []uint8) bool {
	// sequence header only for h264
	if !VideoIsH264(data) {
		return false
	}

	// 2bytes required.
	if len(data) < 2 {
		return false
	}

	frameType := data[0]
	frameType = (frameType >> 4) & 0x0F

	avcPacketType := data[1]

	return avcPacketType == pt.RtmpCodecVideoAVCTypeSequenceHeader
}
