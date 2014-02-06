/*
**
** Author(s):
**  - Cedric GESTES <gestes@aldebaran-robotics.com>
**
** Copyright (C) 2010, 2011, 2012 Aldebaran Robotics
*/

// BUG(r): This is not finished
package messaging

import ("fmt"
	"reflect"
	"bytes")

type QiDecoder interface {
	QiDecode([]byte) error
}

type QiEncoder interface {
	QiEncode() ([]byte, error)
}

type Decoder struct {
	reader *bytes.Buffer
}

type Encoder struct {
	writer *bytes.Buffer
}

func encodeInt(value int64, sz int) []byte {
	b := make([]byte, sz)
	for i := 0 ; i < sz; i++ {
		b[i] = uint8(value)
		value >>= 8
	}
	return b
}

func encodeUInt(value uint64, sz int) []byte {
	b := make([]byte, sz)
	for i := 0; i < sz; i++ {
		b[i] = uint8(value)
		value >>= 8
	}
	return b
}

func NewEncoder(b *bytes.Buffer) *Encoder {
	enc := &Encoder{}
	enc.writer = b
	return enc
}

func (e *Encoder) Encode(value interface{}) error {
	return e.EncodeValue(reflect.ValueOf(value))
}

//todo handle all errors
func (e *Encoder) EncodeValue(value reflect.Value) error {
	switch value.Kind() {
	case reflect.Bool:
		var x = value.Interface().(bool)
		var i int64
		if x { i = 1 } else { i = 0 }
		e.writer.Write(encodeInt(int64(i), 1))
		return nil
	case reflect.Int8:
		var x = value.Interface().(int8)
		e.writer.Write(encodeInt(int64(x), 1))
		return nil
	case reflect.Int16:
		var x = value.Interface().(int16)
		e.writer.Write(encodeInt(int64(x), 2))
		return nil
	case reflect.Int32:
		var x = value.Interface().(int32)
		e.writer.Write(encodeInt(int64(x), 4))
		return nil
	case reflect.Int64:
		var x = value.Interface().(int64)
		e.writer.Write(encodeInt(int64(x), 8))
		return nil
	case reflect.Uint8:
		var x = value.Interface().(uint8)
		e.writer.Write(encodeUInt(uint64(x), 1))
		return nil
	case reflect.Uint16:
		var x = value.Interface().(uint16)
		e.writer.Write(encodeUInt(uint64(x), 2))
		return nil
	case reflect.Uint32:
		var x = value.Interface().(uint32)
		e.writer.Write(encodeUInt(uint64(x), 4))
		return nil
	case reflect.Uint64:
		var x = value.Interface().(uint64)
		e.writer.Write(encodeUInt(uint64(x), 8))
		return nil
	case reflect.String:
		var x string = value.Interface().(string)
		e.writer.Write(encodeUInt(uint64(len(x)), 4))
		e.writer.Write([]byte(x))
		return nil
	}
	return fmt.Errorf("No QiEncoder for type: %s", value.String())
}

func NewDecoder(b *bytes.Buffer) *Decoder {
	enc := &Decoder{}
	enc.reader = b
	return enc
}


//TODO: OPTIMISE
func decodeInt(b *bytes.Buffer, size uint) int64 {
	var v int64 = 0
	var i uint
	var buf [8]byte
	for i = 0; i < size; i++ {
		buf[i], _ = b.ReadByte()
	}
	for i = 0; i < size; i++ {
		c := buf[size - i - 1]
		v = v << 8 | int64(c)
	}
	return v
}

//TODO: OPTIMISE
func decodeUInt(b *bytes.Buffer, size uint) uint64 {
	var v uint64 = 0
	var i uint
	var buf [8]byte
	for i = 0; i < size; i++ {
		buf[i], _ = b.ReadByte()
	}
	for i = 0; i < size; i++ {
		c := buf[size - i - 1]
		fmt.Printf("Intb: %d\n", c)
		v = v << 8 | uint64(c)
	}
	return v
}

func (d *Decoder) DecodeValue(value reflect.Value) error {
	elem := value.Elem()
	if value.Elem().CanSet() == false {
		return fmt.Errorf("The value is not settable, you should pass a pointer")
	}
	switch elem.Kind() {
	case reflect.Bool:
		elem.SetBool(decodeInt(d.reader, 1) != 0)
		return nil

	case reflect.Int8:
		elem.SetInt(decodeInt(d.reader, 1))
		return nil
	case reflect.Int16:
		elem.SetInt(decodeInt(d.reader, 2))
		return nil
	case reflect.Int32:
		elem.SetInt(decodeInt(d.reader, 4))
		return nil
	case reflect.Int64:
		elem.SetInt(decodeInt(d.reader, 8))
		return nil

	case reflect.Uint8:
		elem.SetUint(decodeUInt(d.reader, 1))
		return nil
	case reflect.Uint16:
		elem.SetUint(decodeUInt(d.reader, 2))
		return nil
	case reflect.Uint32:
		elem.SetUint(decodeUInt(d.reader, 4))
		return nil
	case reflect.Uint64:
		elem.SetUint(decodeUInt(d.reader, 8))
		return nil

	case reflect.String:
		var sz uint64
		sz = decodeUInt(d.reader, 4)
		sl := make([]byte, sz)
		rsz, _ := d.reader.Read(sl)
		if uint64(rsz) != sz {
			return fmt.Errorf("Size mismatch...")
		}
		elem.SetString(string(sl))
		return nil

	//array, slice, map, struct
	case reflect.Map:
		fmt.Printf("DecodeValueMap\n")
		return nil
	}
	return fmt.Errorf("no QiDecoder for type: %s", elem.Type().String())
}

func (d *Decoder) Decode(i interface{}) error {
	return d.DecodeValue(reflect.ValueOf(i))
}
