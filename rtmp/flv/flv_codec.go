package flv

import (
	"seal/rtmp/pt"
)

func AudioIsSequenceHeader(data []uint8) bool {

	if !audioIsAAC(data) {
		return false
	}

	if len(data) < 2 {
		return false
	}

	aacPacketType := data[1]

	return aacPacketType == pt.SrsCodecAudioTypeSequenceHeader
}

func audioIsAAC(data []uint8) bool {

	if len(data) < 1 {
		return false
	}

	soundFormat := data[0]
	soundFormat = (soundFormat >> 4) & 0x0f

	return soundFormat == pt.SrsCodecAudioAAC
}

func VideoIsH264(data []uint8) bool {

	if len(data) < 1 {
		return false
	}

	codecId := data[0]
	codecId &= 0x0f

	return pt.SrsCodecVideoAVC == codecId
}

func VideoH264IsKeyframe(data []uint8) bool {
	// 2bytes required.
	if len(data) < 2 {
		return false
	}

	frameType := data[0]
	frameType = (frameType >> 4) & 0x0F

	return frameType == pt.SrsCodecVideoAVCFrameKeyFrame
}

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

	return frameType == pt.SrsCodecVideoAVCFrameKeyFrame && avcPacketType == pt.SrsCodecVideoAVCTypeSequenceHeader
}

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

	return avcPacketType == pt.SrsCodecVideoAVCTypeSequenceHeader
}
