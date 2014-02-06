/*
**
** Author(s):
**  - Cedric GESTES <gestes@aldebaran-robotics.com>
**
** Copyright (C) 2010, 2011, 2012 Aldebaran Robotics
 */

package qi

/*
#cgo CFLAGS: -I/home/ctaf/src/qi-master/lib/libqi -I/home/ctaf/src/qi-master/lib/libqimessaging/c/
#cgo LDFLAGS: -L/home/ctaf/src/qi-master/lib/libqimessaging/build-sys-linux-x86_64/sdk/lib/ -lqic

#include <qic/type.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"reflect"
)

//GoType -> QiType
//QiType -> GoType

func goType(t *C.qi_type_t) (reflect.Type, error) {
	switch C.qi_type_get_kind(t) {
	case C.QI_TYPE_KIND_INVALID:
		return nil, errors.New("Can't convert an invalid type")
	case C.QI_TYPE_KIND_VOID:
		return reflect.TypeOf(nil), nil
	case C.QI_TYPE_KIND_INT:
		return reflect.TypeOf(int(0)), nil
	case C.QI_TYPE_KIND_FLOAT:
		return reflect.TypeOf(float64(0.0)), nil
	case C.QI_TYPE_KIND_STRING:
		return reflect.TypeOf(string("")), nil
	case C.QI_TYPE_KIND_OBJECT:
		//v := interface
		return reflect.TypeOf(AnyObject(nil)), nil
	case C.QI_TYPE_KIND_POINTER:
		elm := C.qi_type_get_value(t)
		gelm, err := goType(elm)
		if err != nil {
			return nil, err
		}
		return reflect.TypeOf(gelm), nil
	case C.QI_TYPE_KIND_DYNAMIC:
		return reflect.TypeOf(AnyObject(nil)), nil
	case C.QI_TYPE_KIND_LIST:
		elm := C.qi_type_get_value(t)
		gelm, err := goType(elm)
		if err != nil {
			return nil, err
		}
		return reflect.SliceOf(gelm), nil
	case C.QI_TYPE_KIND_MAP:
		key := C.qi_type_get_key(t)
		gkey, err := goType(key)
		if err != nil {
			return nil, err
		}
		elm := C.qi_type_get_value(t)
		gelm, err := goType(elm)
		if err != nil {
			return nil, err
		}
		return reflect.MapOf(gkey, gelm), nil
	case C.QI_TYPE_KIND_TUPLE:
		v := make(chan interface{})
		return reflect.SliceOf(reflect.TypeOf(v)), nil
	default:
		return nil, errors.New("error")
	}
}

func qiType(in reflect.Type) (*C.qi_type_t, error) {
	switch in.Kind() {
	case reflect.Invalid:
		return nil, fmt.Errorf("Unsupported type: %v", in.Kind())
	case reflect.Bool:
		return C.qi_type_of_kind(C.QI_TYPE_KIND_BOOL), nil
	case reflect.Int8:
		return C.qi_type_of_kind(C.QI_TYPE_KIND_INT8), nil
	case reflect.Int16:
		return C.qi_type_of_kind(C.QI_TYPE_KIND_INT16), nil
	case reflect.Int32:
		return C.qi_type_of_kind(C.QI_TYPE_KIND_INT32), nil
	case reflect.Int64, reflect.Int:
		return C.qi_type_of_kind(C.QI_TYPE_KIND_INT64), nil
	case reflect.Uint8:
		return C.qi_type_of_kind(C.QI_TYPE_KIND_UINT8), nil
	case reflect.Uint16:
		return C.qi_type_of_kind(C.QI_TYPE_KIND_UINT16), nil
	case reflect.Uint32:
		return C.qi_type_of_kind(C.QI_TYPE_KIND_UINT32), nil
	case reflect.Uint64, reflect.Uint:
		return C.qi_type_of_kind(C.QI_TYPE_KIND_UINT64), nil
	case reflect.Float32:
		return C.qi_type_of_kind(C.QI_TYPE_KIND_FLOAT32), nil
	case reflect.Float64:
		return C.qi_type_of_kind(C.QI_TYPE_KIND_FLOAT64), nil
	case reflect.Complex64, reflect.Complex128:
		//TODO: implement me
		return nil, fmt.Errorf("Unimplemented type: %v", in)
	case reflect.String:
		return C.qi_type_of_kind(C.QI_TYPE_KIND_STRING), nil
	case reflect.Map:
		k, err := qiType(in.Key())
		if err != nil {
			return nil, fmt.Errorf("Unsupported type: %v", in)
		}
		v, err := qiType(in.Elem())
		if err != nil {
			return nil, fmt.Errorf("Unsupported type: %v", in)
		}
		return C.qi_type_map_of(k, v), nil
	case reflect.Slice, reflect.Array:
		v, err := qiType(in.Elem())
		if err != nil {
			return nil, fmt.Errorf("Unsupported slice element type: %v: %s", in.Elem(), err)
		}
		return C.qi_type_list_of(v), nil
	case reflect.Interface, reflect.Struct:
		//TODO: implement me
		return nil, fmt.Errorf("Unimplemented type: %v", in)
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.UnsafePointer:
		return nil, fmt.Errorf("Unsupported type: %v", in.Kind())
	}
	return nil, fmt.Errorf("Unsupported type: %v", in.Kind())
}

//TODO: remove me
func signature(in reflect.Type) (string, error) {
	switch in.Kind() {
	case reflect.Invalid:
		return "X", fmt.Errorf("Unsupported type: %v", in.Kind())
	case reflect.Bool:
		return "b", nil
	case reflect.Int8:
		return "c", nil
	case reflect.Int16:
		return "w", nil
	case reflect.Int32:
		return "i", nil
	case reflect.Int64, reflect.Int:
		return "l", nil
	case reflect.Uint8:
		return "C", nil
	case reflect.Uint16:
		return "W", nil
	case reflect.Uint32:
		return "I", nil
	case reflect.Uint64, reflect.Uint:
		return "L", nil
	case reflect.Float32:
		return "f", nil
	case reflect.Float64:
		return "d", nil
	case reflect.Complex64, reflect.Complex128:
		return "X", fmt.Errorf("Unsupported type: %v", in.Kind())
	case reflect.String:
		return "s", nil
	case reflect.Map:
		k, err := signature(in.Key())
		if err != nil {
			return "X", fmt.Errorf("Unsupported type: %v", in.Kind())
		}
		v, err := signature(in.Elem())
		if err != nil {
			return "X", fmt.Errorf("Unsupported type: %v", in.Kind())
		}
		return "{" + k + v + "}", nil
	case reflect.Slice, reflect.Array:
		v, err := signature(in.Elem())
		if err != nil {
			return "X", fmt.Errorf("Unsupported type: %v", in.Kind())
		}
		return "[" + v + "]", nil
	case reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Ptr, reflect.Struct,
		reflect.UnsafePointer:
		return "X", fmt.Errorf("Unsupported type: %v", in.Kind())
	}
	return "X", fmt.Errorf("Unsupported type: %v", in.Kind())
}

func methodSignature(meth reflect.Method) (string, error) {
	methType := meth.Type

	var sigout string
	var err error
	switch methType.NumOut() {
	case 0:
		sigout = "v"
	case 1:
		tout := methType.Out(0)
		if sigout, err = signature(tout); err != nil {
			return "", err
		}
	default:
		sigout = "("
		for j := 0; j < methType.NumOut(); j++ {
			tout := methType.Out(j)
			fmt.Println("converting(o):", tout)
			sig, err := signature(tout)
			if err != nil {
				return "", err
			}
			sigout += sig
		}
		sigout += ")"
	}

	//TODO: handle method with varargs  type.IsVariadic()
	sigin := "("
	for j := 1; j < methType.NumIn(); j++ {
		tin := methType.In(j)
		fmt.Println("converting(i):", tin)
		sig, err := signature(tin)
		if err != nil {
			return "", err
		}
		sigin += sig
	}
	sigin += ")"
	sig := meth.Name + "::" + sigout + sigin
	return sig, nil
}
