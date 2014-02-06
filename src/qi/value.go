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

#include <qic/value.h>
#include <qic/type.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

//### Decoder

func toGoList(v *C.qi_value_t) (interface{}, error) {
	l := int(C.qi_value_list_size(v))
	gotype, err := goType(C.qi_value_get_type(v))
	if err != nil {
		return nil, err
	}
	sli := reflect.MakeSlice(gotype, l, l)
	for i := 0; i < l; i++ {
		elm := C.qi_value_list_get(v, C.uint(i))
		gelm, err := toGoValue(elm)
		if err != nil {
			return nil, err
		}
		sli.Index(i).Set(reflect.ValueOf(gelm))
	}
	return sli.Interface(), nil
}

func toGoMap(v *C.qi_value_t) (interface{}, error) {
	ckeys := C.qi_value_map_keys(v)
	l := int(C.qi_value_list_size(ckeys))

	gotype, err := goType(C.qi_value_get_type(v))
	if err != nil {
		return nil, err
	}
	ma := reflect.MakeMap(gotype)
	for i := 0; i < l; i++ {
		key := C.qi_value_list_get(ckeys, C.uint(i))
		gkey, err := toGoValue(key)
		if err != nil {
			return nil, err
		}
		val := C.qi_value_map_get(v, key)
		gval, err := toGoValue(val)
		if err != nil {
			return nil, err
		}
		ma.SetMapIndex(reflect.ValueOf(gkey), reflect.ValueOf(gval))
	}
	return ma.Interface(), nil
}

func toGoFromTuple(v *C.qi_value_t) (interface{}, error) {
	l := int(C.qi_value_tuple_size(v))
	sli := make([]interface{}, l)
	for i := 0; i < l; i++ {
		elm := C.qi_value_tuple_get(v, C.uint(i))
		gelm, err := toGoValue(elm)
		if err != nil {
			return nil, err
		}
		sli[i] = gelm
	}
	return sli, nil
}

func toGoValue(v *C.qi_value_t) (interface{}, error) {
	//runtime.setFinaliser(v, C.qi_value_destroy)
	kind := C.qi_value_get_kind(v)
	log.Printf("toGoValue(kind:%d)\n", kind)
	switch kind {
	case C.QI_VALUE_KIND_INVALID:
		return nil, errors.New("Invalid Value")
	case C.QI_VALUE_KIND_VOID:
		return nil, nil
	case C.QI_VALUE_KIND_INT:
		var i C.longlong
		C.qi_value_get_int64(v, &i)

		t := C.qi_value_get_type(v)
		signed := C.qi_type_is_signed(t)
		bits := C.qi_type_get_bits(t)
		var sval int = int(bits)
		if signed == 1 {
			sval = -int(bits)
		}
		switch sval {
		case -64:
			return int64(i), nil
		case -32:
			return int32(i), nil
		case -16:
			return int16(i), nil
		case -8:
			return int8(i), nil
		case 0:
			return bool(i == 1), nil
		case 8:
			return uint8(i), nil
		case 16:
			return uint16(i), nil
		case 32:
			return uint32(i), nil
		case 64:
			return uint64(i), nil
		default:
			return nil, fmt.Errorf("Invalid Bits/Signedness for int, bits:", bits, ",signed:", signed)
		}

	case C.QI_VALUE_KIND_FLOAT:
		var d C.double
		C.qi_value_get_double(v, &d)
		t := C.qi_value_get_type(v)
		bits := C.qi_type_get_bits(t)
		switch bits {
		case 32:
			return float32(d), nil
		case 64:
			return float64(d), nil
		default:
			return nil, errors.New("Invalid Bits for float")
		}
	case C.QI_VALUE_KIND_STRING:
		return C.GoString(C.qi_value_get_string(v)), nil
	case C.QI_VALUE_KIND_OBJECT:
		return toObject(v)
	case C.QI_VALUE_KIND_POINTER:
		return nil, errors.New("Kind pointer not supported")
	case C.QI_VALUE_KIND_DYNAMIC:
		return toObject(v)
		elm := C.qi_value_dynamic_get(v)
		if elm == nil {
			return nil, errors.New("Cant get the element")
		}
		log.Printf("Extracting elem from dynamic(%d)\n", C.qi_value_get_kind(elm))
		return toGoValue(elm)
	case C.QI_VALUE_KIND_LIST:
		return toGoList(v)
	case C.QI_VALUE_KIND_MAP:
		return toGoMap(v)
	case C.QI_VALUE_KIND_TUPLE:
		return toGoFromTuple(v)
		return nil, errors.New("Tuple not implemented yet")
	}
	log.Printf("Cant convert: kind is not supported %d", kind)
	return nil, fmt.Errorf("Cant convert: kind is not supported %d", kind)
}

//#### Encoder
func newValueFromSignature(signature string) *C.qi_value_t {
	//TODO: defer deletes
	return C.qi_value_create(C.CString(signature))
}

func valueSetInt(v *C.qi_value_t, val int64) (*C.qi_value_t, error) {
	ret := C.qi_value_set_int64(v, C.longlong(val))
	if ret == 0 {
		return nil, errors.New("Cant set int value")
	}
	return v, nil
}

func valueSetUInt(v *C.qi_value_t, val uint64) (*C.qi_value_t, error) {
	ret := C.qi_value_set_uint64(v, C.ulonglong(val))
	if ret == 0 {
		return nil, errors.New("Cant set uint value")
	}
	return v, nil
}

func valueSetDouble(v *C.qi_value_t, val float64) (*C.qi_value_t, error) {
	ret := C.qi_value_set_double(v, C.double(val))
	if ret == 0 {
		return nil, errors.New("Cant set double value")
	}
	return v, nil
}

func valueSetString(v *C.qi_value_t, str string) (*C.qi_value_t, error) {
	ret := C.qi_value_set_string(v, C.CString(str))
	if ret == 0 {
		fmt.Printf("Cant set string to val %s\n", str)
		return nil, errors.New("Cant set string value")
	}
	return v, nil
}

func valueSetList(v *C.qi_value_t, val reflect.Value) (*C.qi_value_t, error) {
	l := val.Len()
	for i := 0; i < l; i++ {
		elem, err := toQiValueFromVal(val.Index(i))
		if err != nil {
			return nil, err
		}
		ret := C.qi_value_list_push_back(v, elem)
		if ret != 1 {
			return nil, errors.New("Cant add value to a list")
		}
	}
	return v, nil
}

func valueSetMap(v *C.qi_value_t, val reflect.Value) (*C.qi_value_t, error) {
	for _, kval := range val.MapKeys() {
		key, err := toQiValueFromVal(kval)
		if err != nil {
			return nil, err
		}
		elem, err := toQiValueFromVal(val.MapIndex(kval))
		if err != nil {
			return nil, err
		}
		ret := C.qi_value_map_set(v, key, elem)
		if ret == 0 {
			return nil, errors.New("Cant set map value")
		}
	}
	return v, nil
}

func valueSetTuple(v *C.qi_value_t, val reflect.Value) (*C.qi_value_t, error) {
	return valueSetList(v, val)
	return nil, errors.New("Cant set tuple value")
}

//set the GoValue to the C Value
func valueSet(ret *C.qi_value_t, value reflect.Value) (*C.qi_value_t, error) {
	switch value.Kind() {
	case reflect.Bool:
		if value.Bool() {
			return valueSetInt(ret, 1)
		} else {
			return valueSetInt(ret, 0)
		}
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return valueSetInt(ret, value.Int())
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return valueSetUInt(ret, value.Uint())
	case reflect.Float32:
		return valueSetDouble(ret, value.Float())
	case reflect.Float64:
		return valueSetDouble(ret, value.Float())
	case reflect.String:
		return valueSetString(ret, value.String())
	case reflect.Slice:
		return valueSetList(ret, value)
	case reflect.Map:
		return valueSetMap(ret, value)
	case reflect.Struct:
		return valueSetTuple(ret, value)
	}
	return nil, fmt.Errorf("No QiEncoder for type: %s", value.String())
}

func toQiValue(v interface{}) (*C.qi_value_t, error) {
	return toQiValueFromVal(reflect.ValueOf(v))
}

func toQiValueFromVal(value reflect.Value) (*C.qi_value_t, error) {
	sig, err := signature(value.Type())
	if err != nil {
		return nil, err
	}
	ret := newValueFromSignature(sig)
	fmt.Println("Converting from", value.Type(), "to", sig)
	return valueSet(ret, value)
}

func createTupleValue(in ...interface{}) (*C.qi_value_t, error) {
	sig := "("
	for _, v := range in {
		elm, err := signature(reflect.TypeOf(v))
		if err != nil {
			return nil, err
		}
		sig += elm
	}
	sig += ")"
	fmt.Printf("Signature: %s\n", sig)
	tupleval := newValueFromSignature(sig)

	for i, v := range in {
		cval, err := toQiValue(v)
		if err != nil {
			fmt.Printf("Value not good enough to be allowed to survice\n")
			return nil, err
		}
		C.qi_value_tuple_set(tupleval, C.uint(i), cval)
	}
	return tupleval, nil
}
