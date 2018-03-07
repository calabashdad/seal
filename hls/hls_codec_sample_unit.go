package hls

// the codec sample unit.
// for h.264 video packet, a NALU is a sample unit.
// for aac raw audio packet, a NALU is the entire aac raw data.
// for sequence header, it's not a sample unit.
type codecSampleUnit struct {
	payloadSize int
	payload     []byte
}
