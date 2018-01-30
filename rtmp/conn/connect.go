package conn

import (
	"UtilsTools/identify_panic"
	"fmt"
	"log"
	"seal/rtmp/pt"
)

func (rc *RtmpConn) Connect() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	var connectPkg *pt.ConnectPacket
	var pkt pt.Packet
	pkt = connectPkg
	err = rc.ExpectMsg(&pkt)
	if err != nil {
		log.Println("expect connect packet failed. err=", err)
		return
	}

	connectPkg = pkt.(*pt.ConnectPacket)

	if nil == connectPkg.GetObjectProperty("tcUrl") {
		err = fmt.Errorf("no tcUrl info in connect.")
		return
	}
	rc.ConnectInfo.TcUrl = connectPkg.GetObjectProperty("tcUrl").(string)
	if o := connectPkg.GetObjectProperty("pageUrl"); o != nil {
		rc.ConnectInfo.PageUrl = o.(string)
	}
	if o := connectPkg.GetObjectProperty("swfUrl"); o != nil {
		rc.ConnectInfo.SwfUrl = o.(string)
	}
	if o := connectPkg.GetObjectProperty("objectEncoding"); o != nil {
		rc.ConnectInfo.ObjectEncoding = o.(float64)
	}

	log.Println("expect connect pkt success.", connectPkg, ", info=", rc.ConnectInfo)

	return
}
