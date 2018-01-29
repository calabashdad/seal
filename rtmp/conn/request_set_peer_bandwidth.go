package conn

import (
	"UtilsTools/identify_panic"
	"log"
	"seal/rtmp/protocol"
)

/**
 * @type: The sender can mark this message hard (0), soft (1), or dynamic (2)
 * using the Limit type field.
 */
func (rc *RtmpConn) SetPeerBandWidth(bandwidth uint32, typeLimit uint8) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	var pkt protocol.SetPeerBandWidthPacket

	pkt.Bandwidth = bandwidth
	pkt.TypeLimit = typeLimit

	err = rc.SendPacket(&pkt, 0)
	if err != nil {
		return
	}

	return
}
