package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

type User struct {
	ID       int
	RealName string `unpack:"-"`
	Login    string
	Flags    int
}

func UnpackReflect(u interface{}, data []byte) error {
	r := bytes.NewReader(data)

	val := reflect.ValueOf(u)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("input is not pointer")
	}

	val = val.Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		if typeField.Tag.Get("unpack") == "-" {
			continue
		}

		switch typeField.Type.Kind() {
		case reflect.Int:
			var value uint32
			binary.Read(r, binary.LittleEndian, &value)
			valueField.SetInt(int64(value))
		case reflect.String:
			var lenRaw uint32
			binary.Read(r, binary.LittleEndian, &lenRaw)

			dataRaw := make([]byte, lenRaw)
			binary.Read(r, binary.LittleEndian, &dataRaw)

			valueField.SetString(string(dataRaw))
		default:
			return fmt.Errorf("bad type: %v for field %v", typeField.Type.Kind(), typeField.Name)
		}
	}

	return nil
}

func main() {
	/*
		perl -E '$b = pack("L L/a* L", 1_123_456, "kek.kekovitch", 16);
			print map { ord.", "  } split("", $b); '
	*/
	data := []byte{
		128, 36, 17, 0,

		13, 0, 0, 0,
		107, 101, 107, 46, 107, 101, 107, 111, 118, 105, 116, 99, 104,

		16, 0, 0, 0,
	}
	u := &User{}
	err := UnpackReflect(u, data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", u)
}
