//Package pt is protocol for short. define the basic rtmp protocol consts
package pt

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
)

type Amf0Object struct {
	PropertyName string
	Value        interface{}
	ValueType    uint8 //just for help known type.
}

type amf0EcmaArray struct {
	count     uint32
	anyObject []Amf0Object
}

type amf0StrictArray struct {
	count     uint32
	anyObject []interface{}
}

func (array *amf0EcmaArray) addObject(obj Amf0Object) {
	array.anyObject = append(array.anyObject, obj)
	array.count++
}

// this function do not affect the offset parsed in data.
func amf0ObjectEOF(data []uint8, offset *uint32) (res bool) {
	if len(data) < 3 {
		res = false
		return
	}

	if 0x00 == data[*offset] &&
		0x00 == data[*offset+1] &&
		RTMP_AMF0_ObjectEnd == data[*offset+2] {
		res = true
		*offset += 3
	} else {
		res = false
	}

	return
}

func amf0ReadUtf8(data []uint8, offset *uint32) (value string, err error) {
	if (uint32(len(data)) - *offset) < 2 {
		err = fmt.Errorf("Amf0ReadUtf8: 1, data len is not enough")
		return
	}

	dataLen := binary.BigEndian.Uint16(data[*offset : *offset+2])
	*offset += 2

	if (uint32(len(data)) - *offset) < uint32(dataLen) {
		err = fmt.Errorf("Amf0ReadUtf8: 2, data len is not enough")
		return
	}

	if 0 == dataLen {
		return
	}

	value = string(data[*offset : *offset+uint32(dataLen)])
	*offset += uint32(dataLen)

	return
}

func amf0WriteUtf8(value string) (data []uint8) {

	data = make([]uint8, 2+len(value))

	var offset uint32

	binary.BigEndian.PutUint16(data[offset:offset+2], uint16(len(value)))
	offset += 2

	copy(data[offset:], value)
	offset += uint32(len(value))

	return
}

func amf0ReadAny(data []uint8, marker *uint8, offset *uint32) (value interface{}, err error) {

	if amf0ObjectEOF(data, offset) {
		return
	}

	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0ReadAny: 0, data len is not enough")
		return
	}

	*marker = data[*offset]

	switch *marker {
	case RTMP_AMF0_String:
		value, err = Amf0ReadString(data, offset)
	case RTMP_AMF0_Boolean:
		value, err = amf0ReadBool(data, offset)
	case RTMP_AMF0_Number:
		value, err = Amf0ReadNumber(data, offset)
	case RTMP_AMF0_Null:
		err = amf0ReadNull(data, offset)
	case RTMP_AMF0_Undefined:
		err = amf0ReadUndefined(data, offset)
	case RTMP_AMF0_Object:
		value, err = amf0ReadObject(data, offset)
	case RTMP_AMF0_LongString:
		value, err = amf0ReadLongString(data, offset)
	case RTMP_AMF0_EcmaArray:
		value, err = amf0ReadEcmaArray(data, offset)
	case RTMP_AMF0_StrictArray:
		value, err = amf0ReadStrictArray(data, offset)
	default:
		err = fmt.Errorf("Amf0ReadAny: unknown marker Value, marker=%d", marker)
	}

	if err != nil {
		return
	}

	return
}

func amf0WriteAny(any Amf0Object) (data []uint8) {
	switch any.ValueType {
	case RTMP_AMF0_String:
		data = amf0WriteString(any.Value.(string))
	case RTMP_AMF0_Boolean:
		data = amf0WriteBool(any.Value.(bool))
	case RTMP_AMF0_Number:
		data = amf0WriteNumber(any.Value.(float64))
	case RTMP_AMF0_Null:
		data = amf0WriteNull()
	case RTMP_AMF0_Undefined:
		data = amf0WriteUndefined()
	case RTMP_AMF0_Object:
		data = amf0WriteObject(any.Value.([]Amf0Object))
	case RTMP_AMF0_LongString:
		data = amf0WriteLongString(any.Value.(string))
	case RTMP_AMF0_EcmaArray:
		data = amf0WriteEcmaArray(any.Value.(amf0EcmaArray))
	case RTMP_AMF0_StrictArray:
		data = amf0WriteStrictArray(any.Value.([]Amf0Object))
	default:
		log.Println("Amf0WriteAny: unknown type.", any.ValueType)
	}
	return
}

func Amf0ReadString(data []uint8, offset *uint32) (value string, err error) {

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

	value, err = amf0ReadUtf8(data, offset)

	return
}

func amf0WriteString(value string) (data []uint8) {

	data = append(data, RTMP_AMF0_String)
	data = append(data, amf0WriteUtf8(value)...)

	return
}

func Amf0ReadNumber(data []uint8, offset *uint32) (value float64, err error) {

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

func amf0WriteNumber(value float64) (data []uint8) {
	data = make([]uint8, 1+8)

	var offset uint32

	data[offset] = RTMP_AMF0_Number
	offset += 1

	v2 := math.Float64bits(value)
	binary.BigEndian.PutUint64(data[offset:offset+8], v2)
	offset += 8

	return
}

func amf0ReadObject(data []uint8, offset *uint32) (amf0objects []Amf0Object, err error) {

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

		if amf0ObjectEOF(data, offset) {
			break
		}

		var amf0object Amf0Object

		amf0object.PropertyName, err = amf0ReadUtf8(data, offset)
		if err != nil {
			break
		}

		amf0object.Value, err = amf0ReadAny(data, &amf0object.ValueType, offset)
		if err != nil {
			break
		}

		amf0objects = append(amf0objects, amf0object)
	}

	return
}

func amf0WriteObject(amf0objects []Amf0Object) (data []uint8) {

	data = append(data, RTMP_AMF0_Object)

	for _, v := range amf0objects {
		data = append(data, amf0WriteUtf8(v.PropertyName)...)
		data = append(data, amf0WriteAny(v)...)
	}

	data = append(data, 0x00, 0x00, 0x09)

	return
}

func amf0ReadBool(data []uint8, offset *uint32) (value bool, err error) {
	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0ReadBool:  0, data len is not enough")
		return
	}

	marker := data[*offset]
	*offset++

	if RTMP_AMF0_Boolean != marker {
		err = fmt.Errorf("Amf0ReadBool: RTMP_AMF0_Boolean != marker")
		return
	}

	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0ReadBool:  1, data len is not enough")
		return
	}

	v := data[*offset]
	*offset++

	if v != 0 {
		value = true
	} else {
		value = false
	}

	return
}

func amf0WriteBool(value bool) (data []uint8) {

	data = make([]uint8, 1+1)
	data[0] = RTMP_AMF0_Boolean
	if value {
		data[1] = 1
	} else {
		data[1] = 0
	}

	return
}

func amf0ReadNull(data []uint8, offset *uint32) (err error) {
	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0ReadNull:  0, data len is not enough")
		return
	}

	marker := data[*offset]
	*offset++

	if RTMP_AMF0_Null != marker {
		err = fmt.Errorf("Amf0ReadNull: RTMP_AMF0_Null != marker")
		return
	}

	return
}

func amf0WriteNull() (data []uint8) {

	data = append(data, RTMP_AMF0_Null)

	return
}

func amf0ReadUndefined(data []uint8, offset *uint32) (err error) {

	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0ReadUndefined:  0, data len is not enough")
		return
	}

	marker := data[*offset]
	*offset++

	if RTMP_AMF0_Undefined != marker {
		err = fmt.Errorf("Amf0ReadUndefined: RTMP_AMF0_Undefined != marker")
		return
	}

	return
}

func amf0WriteUndefined() (data []uint8) {

	data = append(data, RTMP_AMF0_Undefined)

	return
}

func amf0ReadLongString(data []uint8, offset *uint32) (value string, err error) {

	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0ReadLongString: 0, data len is not enough")
		return
	}

	marker := data[*offset]
	*offset++

	if RTMP_AMF0_LongString != marker {
		err = fmt.Errorf("Amf0ReadLongString: RTMP_AMF0_LongString != marker")
		return
	}

	if (uint32(len(data)) - *offset) < 4 {
		err = fmt.Errorf("Amf0ReadLongString: 1, data len is not enough")

		return
	}

	dataLen := binary.BigEndian.Uint32(data[*offset : *offset+4])
	*offset += 4
	if dataLen <= 0 {
		err = fmt.Errorf("Amf0ReadLongString: data len is <= 0, dataLen=%d", dataLen)
		return
	}

	if (uint32(len(data)) - *offset) < dataLen {
		err = fmt.Errorf("Amf0ReadLongString: 2, data len is not enough")
		return
	}

	value = string(data[*offset : *offset+dataLen])
	*offset += dataLen

	return
}

func amf0WriteLongString(value string) (data []uint8) {

	data = make([]uint8, 1+4+len(value))

	var offset uint32

	data[offset] = RTMP_AMF0_LongString
	offset++

	dataLen := len(value)
	binary.BigEndian.PutUint32(data[offset:offset+4], uint32(dataLen))
	offset += 4

	copy(data[offset:], value)
	offset += uint32(dataLen)

	return
}

func amf0ReadEcmaArray(data []uint8, offset *uint32) (value amf0EcmaArray, err error) {

	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0ReadEcmaArray: 0, data len is not enough")
		return
	}

	marker := data[*offset]
	*offset++

	if RTMP_AMF0_EcmaArray != marker {
		err = fmt.Errorf("error: Amf0ReadEcmaArray: RTMP_AMF0_EcmaArray != marker")
		return
	}

	if (uint32(len(data)) - *offset) < 4 {
		err = fmt.Errorf("Amf0ReadEcmaArray: 1, data len is not enough")
		return
	}

	value.count = binary.BigEndian.Uint32(data[*offset : *offset+4])
	*offset += 4

	for {
		if *offset >= uint32(len(data)) {
			break
		}

		if amf0ObjectEOF(data, offset) {
			break
		}

		var amf Amf0Object
		amf.PropertyName, err = amf0ReadUtf8(data, offset)
		if err != nil {
			break
		}

		amf.Value, err = amf0ReadAny(data, &amf.ValueType, offset)
		if err != nil {
			break
		}

		value.anyObject = append(value.anyObject, amf)
	}

	return
}

func amf0WriteEcmaArray(arr amf0EcmaArray) (data []uint8) {
	data = make([]uint8, 1+4)

	var offset uint32

	data[offset] = RTMP_AMF0_EcmaArray
	offset++

	binary.BigEndian.PutUint32(data[offset:offset+4], uint32(arr.count))
	offset += 4

	for _, v := range arr.anyObject {
		data = append(data, amf0WriteUtf8(v.PropertyName)...)
		data = append(data, amf0WriteAny(v)...)
	}

	//eof
	data = append(data, 0x00, 0x00, 0x09)

	return
}

func amf0ReadStrictArray(data []uint8, offset *uint32) (value amf0StrictArray, err error) {
	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0ReadStrictArray: 0, data len is not enough")
		return
	}

	marker := data[*offset]
	*offset++

	if RTMP_AMF0_StrictArray != marker {
		err = fmt.Errorf("Amf0ReadStrictArray: error: RTMP_AMF0_StrictArray != marker")
		return
	}

	if (uint32(len(data)) - *offset) < 4 {
		err = fmt.Errorf("Amf0ReadStrictArray: 1, data len is not enough")
		return
	}

	value.count = binary.BigEndian.Uint32(data[*offset : *offset+4])
	*offset += 4

	for i := 0; uint32(i) < value.count; i++ {
		if *offset >= uint32(len(data)) {
			break
		}

		var obj interface{}

		var markerLocal uint8
		obj, err = amf0ReadAny(data, &markerLocal, offset)
		if err != nil {
			break
		}

		value.anyObject = append(value.anyObject, obj)
	}

	return
}

func amf0WriteStrictArray(objs []Amf0Object) (data []uint8) {
	data = make([]uint8, 1+4)

	var offset uint32

	data[offset] = RTMP_AMF0_StrictArray
	offset++

	count := len(objs)
	binary.BigEndian.PutUint32(data[offset:offset+4], uint32(count))
	offset += 4

	for _, v := range objs {
		data = append(data, amf0WriteAny(v)...)
	}

	//eof
	data = append(data, 0x00, 0x00, 0x09)

	return
}

func amf0Discovery(data []uint8, offset *uint32) (value interface{}, marker uint8, err error) {

	if amf0ObjectEOF(data, offset) {
		return
	}

	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0Discovery: 0, data len is not enough")
		return
	}

	marker = data[*offset]

	switch marker {
	case RTMP_AMF0_String:
		value, err = Amf0ReadString(data, offset)
	case RTMP_AMF0_Boolean:
		value, err = amf0ReadBool(data, offset)
	case RTMP_AMF0_Number:
		value, err = Amf0ReadNumber(data, offset)
	case RTMP_AMF0_Null:
		err = amf0ReadNull(data, offset)
	case RTMP_AMF0_Undefined:
		err = amf0ReadUndefined(data, offset)
	case RTMP_AMF0_Object:
		value, err = amf0ReadObject(data, offset)
	case RTMP_AMF0_LongString:
		value, err = amf0ReadLongString(data, offset)
	case RTMP_AMF0_EcmaArray:
		value, err = amf0ReadEcmaArray(data, offset)
	case RTMP_AMF0_StrictArray:
		value, err = amf0ReadStrictArray(data, offset)
	default:
		err = fmt.Errorf("Amf0Discovery: unknown marker type, marker=%d", marker)
	}

	if err != nil {
		return
	}

	return
}
