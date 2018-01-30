package conn

import (
	"UtilsTools/identify_panic"
	"fmt"
	"log"
	"seal/conf"
)

func (rc *RtmpConn) Loop() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	var err error

	err = rc.SetWindowAckSize(2500000)
	if err != nil {
		log.Println("set window ack size error.", err)
		return
	}
	log.Println("send set window ack size success.")

	err = rc.SetPeerBandWidth(2500000, 2)
	if err != nil {
		log.Println("set peer band width error.", err)
		return
	}
	log.Println("send set peer bandwidth success.")

	//todo. bandwidth test.

	err = rc.ResponseConnect()
	if err != nil {
		return
	}
	log.Println("response connect success.")

	//todo.
	// err = rc.OnBwDone()
	// if err != nil {
	// 	return
	// }
	// log.Println("send on bw done success.")

	err = rc.IdentifyClient()
	if err != nil {
		log.Println("identify client type failed.err=", err)
		return
	}

	log.Println("client identify success. role=", rc.Role, ",streamName=", rc.StreamName, ",tokenStr=", rc.TokenStr)

	chunkSize := conf.GlobalConfInfo.Rtmp.ChunkSize
	err = rc.RequestSetChunkSize(chunkSize)
	if err != nil {
		return
	}
	log.Println("send set chunk size success. chunk size")

	switch rc.Role {
	case RtmpRolePlayer:
		err = rc.DoPlayerCycle()
	case RtmpRoleFMLEPublisher:
		err = rc.DoFmlePublisherCycle()
	case RtmpRoleFlashPublisher:
		err = rc.DoFlashPublisherCycle()
	default:
		err = fmt.Errorf("unknown rtmp role.")
		return
	}

	if err != nil {
		return
	}

	return
}
