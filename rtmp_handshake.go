package main

import (
	"UtilsTools/identify_panic"
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"time"
)

var (
	handshakeClientFullKey = []uint8{
		'G', 'e', 'n', 'u', 'i', 'n', 'e', ' ', 'A', 'd', 'o', 'b', 'e', ' ',
		'F', 'l', 'a', 's', 'h', ' ', 'P', 'l', 'a', 'y', 'e', 'r', ' ',
		'0', '0', '1',
		0xF0, 0xEE, 0xC2, 0x4A, 0x80, 0x68, 0xBE, 0xE8, 0x2E, 0x00, 0xD0, 0xD1,
		0x02, 0x9E, 0x7E, 0x57, 0x6E, 0xEC, 0x5D, 0x2D, 0x29, 0x80, 0x6F, 0xAB,
		0x93, 0xB8, 0xE6, 0x36, 0xCF, 0xEB, 0x31, 0xAE,
	}
	handshakeServerFullKey = []uint8{
		'G', 'e', 'n', 'u', 'i', 'n', 'e', ' ', 'A', 'd', 'o', 'b', 'e', ' ',
		'F', 'l', 'a', 's', 'h', ' ', 'M', 'e', 'd', 'i', 'a', ' ',
		'S', 'e', 'r', 'v', 'e', 'r', ' ',
		'0', '0', '1',
		0xF0, 0xEE, 0xC2, 0x4A, 0x80, 0x68, 0xBE, 0xE8, 0x2E, 0x00, 0xD0, 0xD1,
		0x02, 0x9E, 0x7E, 0x57, 0x6E, 0xEC, 0x5D, 0x2D, 0x29, 0x80, 0x6F, 0xAB,
		0x93, 0xB8, 0xE6, 0x36, 0xCF, 0xEB, 0x31, 0xAE,
	}

	handshakeClientPartialKey = handshakeClientFullKey[:30]
	handshakeServerPartialKey = handshakeServerFullKey[:36]
)

func (rtmp *RtmpSession) HandShake() (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
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

	err = rtmp.Conn.SetDeadline(time.Now().Add(time.Duration(g_conf_info.Rtmp.TimeOut) * time.Second))
	if err != nil {
		return
	}

	//recv c0c1
	_, err = io.ReadFull(rtmp.Conn, c0c1)
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
		if !ComplexHandShake(c1, s0, s1, s2) {
			err = fmt.Errorf("0 != clientVer, complex handshake failed.")
			return
		}
	} else {
		//use simple handshake
		log.Println("0 == clientVer, client use simple handshake.")
		s1[0] = 3
		copy(s1, c2)
		copy(s2, c1)
	}

	//send s0s1s2
	err = rtmp.Conn.SetDeadline(time.Now().Add(time.Duration(g_conf_info.Rtmp.TimeOut) * time.Second))
	if err != nil {
		return
	}
	if _, err = rtmp.Conn.Write(s0s1s2); err != nil {
		return
	}

	//recv c2
	err = rtmp.Conn.SetDeadline(time.Now().Add(time.Duration(g_conf_info.Rtmp.TimeOut) * time.Second))
	if err != nil {
		return
	}
	if _, err = io.ReadFull(rtmp.Conn, c2); err != nil {
		return
	}

	//c2 do not need verify.

	return
}

func ComplexHandShake(c1 []uint8, s0 []uint8, s1 []uint8, s2 []uint8) bool {

	clientTime := binary.BigEndian.Uint32(c1[0:4])
	clientVer := binary.BigEndian.Uint32(c1[4:8])

	_ = clientVer

	//use digest-key scheme.
	c1Digest764 := c1[8 : 8+764]
	var serverDigestForS2 []uint8
	if ok, digest := IsDigestKeyScheme(c1, c1Digest764); !ok {
		//failed try key-digest scheme
		c1Digest764_2 := c1[8+764 : 8+764+764]
		if ok2, digest2 := IsKeyDigestScheme(c1, c1Digest764_2); !ok2 {
			log.Println("ComplexHandShake verify both digest-key scheme and key-digest failed.")
			return false
		} else {
			serverDigestForS2 = digest2
		}
	} else {
		serverDigestForS2 = digest
	}

	//create s0
	s0[0] = 3

	//create s1
	serverTime := clientTime
	serverVer := uint32(0x0a0b0c0d)
	binary.BigEndian.PutUint32(s1[0:4], serverTime)
	binary.BigEndian.PutUint32(s1[4:8], serverVer)
	for {
		rand.Read(s1[8:])
		if uint32(s1[8]+s1[9]+s1[10]+s1[11]) < (764 - 32) {
			break
		}
	}

	//use digest-key scheme.
	digestLoc := 8 + s1[8] + s1[9] + s1[10] + s1[11]

	h := hmac.New(sha256.New, handshakeServerPartialKey)
	h.Write(s1[:digestLoc])
	h.Write(s1[digestLoc+32:])
	digestData := h.Sum(nil)
	copy(s1[digestLoc:], digestData)

	//create s2.
	// 1536bytes c2s2. c2 and s2 has the same structure.
	//random-data: 1504bytes
	//digest-data: 32bytes
	rand.Read(s2[:])
	h = hmac.New(sha256.New, serverDigestForS2)
	h.Write(s2[:len(s2)-32])
	digestS2 := h.Sum(nil)
	copy(s2[len(s2)-32:], digestS2)

	return true
}

//just for c1 or s1
func IsDigestKeyScheme(buf []uint8, c1Digest764 []uint8) (ok bool, digest []uint8) {

	// 764bytes digest
	//offset: 4bytes (u[0] + u[1] + u[2] + u[3])
	//random-data: (offset)bytes
	//digest-data: 32bytes
	//random-data: (764-4-offset-32)bytes

	// 764bytes key
	//random-data: (offset)bytes
	//key-data: 128bytes
	//random-data: (764-offset-128-4)bytes
	//offset: 4bytes

	var digestOffset uint32
	for i := 0; i < 4; i++ {
		digestOffset += uint32(c1Digest764[i])
	}

	if digestOffset > (764 - 32) {
		ok = false
		return
	}

	digestLoc := 4 + digestOffset
	digestData := c1Digest764[digestLoc : digestLoc+32]

	//part1 and part2 is divided by digest data of c1 or s1.
	part1 := buf[:8+digestLoc]
	part2 := buf[8+digestLoc+32:]

	h := hmac.New(sha256.New, handshakeClientPartialKey)
	h.Write(part1)
	h.Write(part2)
	calcDigestData := h.Sum(nil)

	if 0 == bytes.Compare(digestData, calcDigestData) {
		ok = true
		h := hmac.New(sha256.New, handshakeServerFullKey)
		h.Write(digestData)
		digest = h.Sum(nil)
	} else {
		ok = false
	}

	return
}

func IsKeyDigestScheme(buf []uint8, c1Digest764 []uint8) (ok bool, digest []uint8) {
	// 764bytes key
	//random-data: (offset)bytes
	//key-data: 128bytes
	//random-data: (764-offset-128-4)bytes
	//offset: 4bytes

	// 764bytes digest
	//offset: 4bytes (u[0] + u[1] + u[2] + u[3])
	//random-data: (offset)bytes
	//digest-data: 32bytes
	//random-data: (764-4-offset-32)bytes

	var digestOffset uint32
	for i := 0; i < 4; i++ {
		digestOffset += uint32(c1Digest764[i])
	}

	if digestOffset > (764 - 32) {
		ok = false
		return
	}

	digestLoc := 4 + digestOffset
	digestData := c1Digest764[digestLoc : digestLoc+32]

	//part1 and part2 is divided by digest data of c1 or s1.
	part1 := buf[:8+764+digestLoc]
	part2 := buf[8+764+digestLoc+32:]

	h := hmac.New(sha256.New, handshakeClientPartialKey)
	h.Write(part1)
	h.Write(part2)
	calcDigestData := h.Sum(nil)

	if 0 == bytes.Compare(digestData, calcDigestData) {
		ok = true
		h := hmac.New(sha256.New, handshakeServerFullKey)
		h.Write(digestData)
		digest = h.Sum(nil)
	} else {
		ok = false
	}

	return
}
