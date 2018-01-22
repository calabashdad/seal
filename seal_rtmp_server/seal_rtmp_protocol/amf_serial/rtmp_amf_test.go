package amf_serial

import "testing"

func TestAmf0WriteUtf8(t *testing.T) {

	s := "abcd"

	data := Amf0WriteUtf8(s)

	var offset uint32
	_, value := Amf0ReadUtf8(data, &offset)

	if string(value) != s {
		t.Error("TestAmf0WriteUtf8 failed.")
	}
}

func TestAmf0WriteAny(t *testing.T) {
	var a Amf0Object

	a.PropertyName = "test"
	a.Value = float64(123)
	a.ValueType = RTMP_AMF0_Number

	data := Amf0WriteAny(a)
	_ = data

	var offset uint32
	_, res := Amf0ReadAny(data, &offset)
	_ = res

	v := res.(float64)

	if v != 123 {
		t.Error("TestAmf0WriteAny 1")
	}

}

func TestAmf0WriteString(t *testing.T) {
	s := "abcd"

	data := Amf0WriteString(s)

	var offset uint32
	_, value := Amf0ReadString(data, &offset)

	if string(value) != s {
		t.Error("TestAmf0WriteString failed.")
	}
}

func TestAmf0WriteNumber(t *testing.T) {
	data := Amf0WriteNumber(123.0)

	var offset uint32
	_, res := Amf0ReadNumber(data, &offset)

	if res != 123.0 {
		t.Error("TestAmf0WriteNumber 0.")
	}
}

func TestAmf0WriteObject(t *testing.T) {
	var objs []Amf0Object

	var obj1 Amf0Object
	obj1.PropertyName = "o1"
	obj1.Value = "obj1 Value"
	obj1.ValueType = RTMP_AMF0_String

	objs = append(objs, obj1)

	var obj2 Amf0Object
	obj2.PropertyName = "o2"
	obj2.Value = 123.0
	obj2.ValueType = RTMP_AMF0_Number

	objs = append(objs, obj2)

	data := Amf0WriteObject(objs)

	var offset uint32
	err, res := Amf0ReadObject(data, &offset)
	if err != nil {
		t.Error("TestAmf0WriteObject ", err)
	}
	_ = res
}

func TestAmf0WriteBool(t *testing.T) {
	data := Amf0WriteBool(false)

	var offset uint32
	_, res := Amf0ReadBool(data, &offset)

	if res {
		t.Error("TestAmf0WriteBool 0.")
	}
}

func TestAmf0WriteNull(t *testing.T) {
	data := Amf0WriteNull()

	var offset uint32
	err := Amf0ReadNull(data, &offset)

	if err != nil {
		t.Error("TestAmf0WriteNull 0.")
	}
}

func TestAmf0WriteUndefined(t *testing.T) {
	data := Amf0WriteUndefined()

	var offset uint32
	err := Amf0ReadUndefined(data, &offset)

	if err != nil {
		t.Error("TestAmf0WriteUndefined 0.")
	}
}

func TestAmf0WriteLongString(t *testing.T) {
	s := "abcd"

	data := Amf0WriteString(s)

	var offset uint32
	_, value := Amf0ReadString(data, &offset)

	if string(value) != s {
		t.Error("TestAmf0WriteString failed.")
	}
}

func TestAmf0WriteEcmaArray(t *testing.T) {
	var objs []Amf0Object

	var obj1 Amf0Object
	obj1.PropertyName = "o1"
	obj1.Value = "obj1 Value"
	obj1.ValueType = RTMP_AMF0_String

	objs = append(objs, obj1)

	var obj2 Amf0Object
	obj2.PropertyName = "o2"
	obj2.Value = 123.0
	obj2.ValueType = RTMP_AMF0_Number

	objs = append(objs, obj2)

	data := Amf0WriteEcmaArray(objs)

	var offfset uint32
	err, _ := Amf0ReadEcmaArray(data, &offfset)
	if err != nil {
		t.Error("TestAmf0WriteEcmaArray.")
	}
}

func TestAmf0WriteStrictArray(t *testing.T) {
	var objs []Amf0Object

	var obj1 Amf0Object
	obj1.PropertyName = "o1"
	obj1.Value = "obj1 Value"
	obj1.ValueType = RTMP_AMF0_String

	objs = append(objs, obj1)

	var obj2 Amf0Object
	obj2.PropertyName = "o2"
	obj2.Value = 123.0
	obj2.ValueType = RTMP_AMF0_Number

	objs = append(objs, obj2)

	data := Amf0WriteStrictArray(objs)

	var offfset uint32
	err, _ := Amf0ReadStrictArray(data, &offfset)
	if err != nil {
		t.Error("TestAmf0WriteStrictArray.")
	}

}

func TestAmf0Discovery(t *testing.T) {

}
