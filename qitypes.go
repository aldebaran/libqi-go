package main

import (
    "fmt"
    "reflect"
)

// const (
//     Invalid Kind = iota
//     Bool
//     Int
//     Int8
//     Int16
//     Int32
//     Int64
//     Uint
//     Uint8
//     Uint16
//     Uint32
//     Uint64
//     Uintptr
//     Float32
//     Float64
//     Complex64
//     Complex128
//     Array
//     Chan
//     Func
//     Interface
//     Map
//     Ptr
//     Slice
//     String
//     Struct
//     UnsafePointer
// )

//func TakeSignature(i interface{}) string {
func TakeSignature(t reflect.Type) string {
	//take the signature of a type
	switch t.Kind() {
	case reflect.Bool:
		return "b";

	case reflect.Int8:
		return "c";
	case reflect.Int16:
		return "w";
	case reflect.Int:
		return "i";
	case reflect.Int32:
		return "i";
	case reflect.Int64:
		return "l";

	case reflect.Uint8:
		return "C";
	case reflect.Uint16:
		return "W";
	case reflect.Uint:
		return "I";
	case reflect.Uint32:
		return "I";
	case reflect.Uint64:
		return "L";

	case reflect.Float32:
		return "f";
	case reflect.Float64:
		return "d";

	case reflect.String:
		return "s";

	case reflect.Array:
		var sig string = "["
		sig += TakeSignature(t.Elem())
		sig = sig + "]"
		return sig

	case reflect.Slice:
		var sig string = "["
		sig += TakeSignature(t.Elem())
		sig = sig + "]"
		return sig

	case reflect.Map:
		var sig string = "{"
		sig += TakeSignature(t.Key())
		sig += TakeSignature(t.Elem())
		sig = sig + "}"
		return sig

	case reflect.Struct:
		var sig string = "("
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			sig += TakeSignature(f.Type)
		}
		sig += ")"
		return sig

// case reflect.UnsafePointer:
// case reflect.Ptr:
// case reflect.Interface:
// case reflect.Chan:
// case reflect.Func:
// case reflect.Complex64:
// case reflect.Complex128:
// case reflect.Uintptr:
	default:
		return "X"
	}
	return "X";
}


type T struct {
	A int
	B string
}

func main() {
	var i8 int8 = 2
	var i16 int16 = 2
	var i32 int32 = 2
	var i64 int64 = 2

	var ui8	 uint8 = 2
	var ui16 uint16 = 2
	var ui32 uint32 = 2
	var ui64 uint64 = 2

	var x32 float32 = 3.4
	var x64 float64 = 3.4
	m := make(map[string]int)
	t := T{23, "skidoo"}

	fmt.Println("Signature:", TakeSignature(reflect.TypeOf(i)))
	fmt.Println("Signature:", TakeSignature(reflect.TypeOf(i6)))
	fmt.Println("Signature:", TakeSignature(reflect.TypeOf(i2)))
	fmt.Println("Signature:", TakeSignature(reflect.TypeOf(i4)))

	fmt.Println("Signature:", TakeSignature(reflect.TypeOf(u8)))
	fmt.Println("Signature:", TakeSignature(reflect.TypeOf(u16)))
	fmt.Println("Signature:", TakeSignature(reflect.TypeOf(u32)))
	fmt.Println("Signature:", TakeSignature(reflect.TypeOf(u64)))

	fmt.Println("Signature:", TakeSignature(reflect.TypeOf(x32)))
	fmt.Println("Signature:", TakeSignature(reflect.TypeOf(x64)))
	fmt.Println("Signature:", TakeSignature(reflect.TypeOf(m)))
	fmt.Println("Signature:", TakeSignature(reflect.TypeOf(t)))
}
