/*
**
** Author(s):
**  - Cedric GESTES <gestes@aldebaran-robotics.com>
**
** Copyright (C) 2010, 2011, 2012 Aldebaran Robotics
*/

// BUG(r): This is not finished
package messaging

import (
    "reflect"
)

//qi.Signature(map[string]int, int8, float32)

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

func TakeSignature(i interface{}) string {
	return TakeSignatureType(reflect.TypeOf(i))
}

func TakeSignatureType(t reflect.Type) string {
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
