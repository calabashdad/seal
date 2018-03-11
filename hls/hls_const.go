package hls

const (
	hlsMaxCodecSample     = 128
	hlsAacSampleRateUnset = 15
)

const (
	// in ms, for HLS aac sync time.
	hlsConfDefaultAacSync = 100
)

// the mpegts header specifed the video/audio pid.
const (
	tsVideoPid = 256
	tsAudioPid = 257
)

// ts aac stream id.
const tsAudioAac = 0xc0

// ts avc stream id.
const tsVideoAvc = 0xe0

// the public data, event HLS disable, others can use it.
// 0 = 5.5 kHz = 5512 Hz
// 1 = 11 kHz = 11025 Hz
// 2 = 22 kHz = 22050 Hz
// 3 = 44 kHz = 44100 Hz
var flvSampleRates = []int{5512, 11025, 22050, 44100}

// the sample rates in the codec,
// in the sequence header.
var aacSampleRates = []int{
	96000, 88200, 64000, 48000,
	44100, 32000, 24000, 22050,
	16000, 12000, 11025, 8000,
	7350, 0, 0, 0}

// @see: ngx_rtmp_hls_audio
// We assume here AAC frame size is 1024
// Need to handle AAC frames with frame size of 960 */
const hlsAacSampleSize = 1024

// max PES packets size to flush the video.
const hlsAudioCacheSize = 1024 * 1024

// @see: NGX_RTMP_HLS_DELAY,
// 63000: 700ms, ts_tbn=90000
// 72000: 800ms, ts_tbn=90000
const hlsAutoDelay = 72000

// drop the segment when duration of ts too small.
const hlsSegmentMinDurationMs = 100

// in ms, for HLS aac flush the audio
const hlsAacDelay = 100
