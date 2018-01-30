package conn

import (
	"UtilsTools/identify_panic"
	"encoding/binary"
	"fmt"
	"log"
	"seal/rtmp/pt"
)

func (rc *RtmpConn) HandShake() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
		}
	}()

	var handshakeData [6146]uint8 // c0(1) + c1(1536) + c2(1536) + s0(1) + s1(1536) + s2(1536)

	c0 := handshakeData[:1]
	c1 := handshakeData[1:1537]
	c2 := handshakeData[1537:3073]

	s0 := handshakeData[3073:3074]
	s1 := handshakeData[3074:4610]
	s2 := handshakeData[4610:6146]

	c0c1 := handshakeData[0:1537]
	s0s1s2 := handshakeData[3073:6146]

	//recv c0c1
	err = rc.TcpConn.ExpectBytesFull(c0c1, 1537)
	if err != nil {
		return
	}

	//parse c0
	if c0[0] != 3 {
		err = fmt.Errorf("client c0 is not 3.")
		return
	}

	//use complex handshake, if complex handshake failed, try use simple handshake
	//parse c1
	clientVer := binary.BigEndian.Uint32(c1[4:8])
	if 0 != clientVer {
		err = pt.ComplexHandShake(c1, s0, s1, s2)
		if err != nil {
			return
		}
		log.Println("complex handshake success.")

	} else {
		//use simple handshake
		log.Println("0 == clientVer, client use simple handshake.")
		s0[0] = 3
		copy(s1, c2)
		copy(s2, c1)
	}

	//send s0s1s2
	err = rc.TcpConn.SendBytes(s0s1s2)
	if err != nil {
		return
	}

	//recv c2
	err = rc.TcpConn.ExpectBytesFull(c2, uint32(len(c2)))
	if err != nil {
		return
	}

	//c2 do not need verify.

	return
}
