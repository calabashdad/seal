package handshake

import (
	"UtilsTools/identify_panic"
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"log"
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

func ComplexHandShake(c1 []uint8, s0 []uint8, s1 []uint8, s2 []uint8) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "-", identify_panic.IdentifyPanic())
		}
	}()

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

	//create s1. time(4B) version(4B) [digest]{random} [key]{random}
	serverTime := clientTime
	serverVer := uint32(0x0d0e0a0d)
	binary.BigEndian.PutUint32(s1[0:4], serverTime)
	binary.BigEndian.PutUint32(s1[4:8], serverVer)
	//use digest-key scheme.

	var randomDataoffset uint32
	for {
		rand.Read(s1[4+4:]) // time(4B)server version(4B)
		randomDataoffset = uint32(s1[8] + s1[9] + s1[10] + s1[11])
		if randomDataoffset > 0 && randomDataoffset < 728 {
			break
		}
	}

	digestLoc := 4 + 4 + 4 + randomDataoffset //time(4B) version(4B) + digest[offset(4B) + random1(offset B) + digest + random2] + key[]

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
