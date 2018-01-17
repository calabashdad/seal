package main

import (
	"encoding/binary"
	"fmt"
	"math"
)

type Amf0Object struct {
	propertyName string
	value        interface{}
}

//this function do not affect the offset parsed in data.
func Amf0ObjectEof(data []uint8) (res bool) {
	if len(data) < 3 {
		res = false
		return
	}

	objectEofFlag := uint32(data[0])<<16 + uint32(data[1])<<8 + uint32(data[2])
	if 0x09 == objectEofFlag {
		res = true
	}

	return
}

func Amf0ReadUtf8(data []uint8, offset *uint32) (err error, value string) {
	if (uint32(len(data)) - *offset) < 2 {
		err = fmt.Errorf("Amf0ReadString: 1, data len is not enough")
		return
	}

	dataLen := binary.BigEndian.Uint16(data[*offset : *offset+2])
	*offset += 2

	if dataLen <= 0 {
		err = fmt.Errorf("Amf0ReadString: dataLen <= 0 ")
		return
	}

	if (uint32(len(data)) - *offset) < uint32(dataLen) {
		err = fmt.Errorf("Amf0ReadString: 2, data len is not enough")
		return
	}

	value = string(data[*offset : *offset+uint32(dataLen)])
	*offset += uint32(dataLen)

	return
}

func Amf0ReadAny(data []uint8, offset *uint32) (err error, value interface{}) {

	if Amf0ObjectEof(data[*offset : *offset+3]) {
		return
	}

	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0ReadAny: 0, data len is not enough")
		return
	}

	marker := data[*offset]

	switch marker {
	case RTMP_AMF0_String:
		err, value = Amf0ReadString(data, offset)
	case RTMP_AMF0_Boolean:
	case RTMP_AMF0_Number:
	case RTMP_AMF0_Null:
	case RTMP_AMF0_Undefined:
	case RTMP_AMF0_Object:
	case RTMP_AMF0_LongString:
	case RTMP_AMF0_EcmaArray:
	case RTMP_AMF0_StrictArray:
	case RTMP_AMF0_Invalid:
	default:
		err = fmt.Errorf("Amf0ReadAny: unknown marker value, marker=", marker)
	}

	if err != nil {
		return
	}

	return
}

func Amf0ReadString(data []uint8, offset *uint32) (err error, value string) {

	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0ReadString: 0, data len is not enough")
		return
	}

	marker := data[*offset]
	*offset += 1

	if RTMP_AMF0_String != marker {
		err = fmt.Errorf("Amf0ReadString: RTMP_AMF0_String != marker")
		return
	}

	err, value = Amf0ReadUtf8(data, offset)

	return
}

func Amf0ReadNumber(data []uint8, offset *uint32) (err error, value float64) {

	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0ReadNumber: 0, data len is not enough")
		return
	}

	marker := data[*offset]
	*offset += 1

	if RTMP_AMF0_Number != marker {
		err = fmt.Errorf("Amf0ReadNumber: RTMP_AMF0_Number != marker")
		return
	}

	if (uint32(len(data)) - *offset) < 8 {
		err = fmt.Errorf("Amf0ReadNumber: 1, data len is not enough")
		return
	}

	value_tmp := binary.BigEndian.Uint64(data[*offset : *offset+8])
	*offset += 8

	value = math.Float64frombits(value_tmp)

	return
}

func Amf0ReadObject(data []uint8, offset *uint32) (err error, amf0objects []Amf0Object) {

	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0ReadObject: 0, data len is not enough")
		return
	}

	marker := data[*offset]
	*offset += 1

	if RTMP_AMF0_Object != marker {
		err = fmt.Errorf("error: Amf0ReadObject:RTMP_AMF0_Object != marker")
		return
	}

	for {
		if *offset >= uint32(len(data)) {
			break
		}

		if Amf0ObjectEof(data[*offset : *offset+3]) {
			break
		}

		var amf0object Amf0Object

		err, amf0object.propertyName = Amf0ReadUtf8(data, offset)
		if err != nil {
			break
		}

		err, amf0object.value = Amf0ReadAny(data, offset)
		if err != nil {
			break
		}

		amf0objects = append(amf0objects, amf0object)
	}

	return
}
