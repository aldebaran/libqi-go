// value_test.go
package qi

import (
	"fmt"
	"reflect"
	"testing"
)

func compareValue(t *testing.T, v, v2 interface{}) {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Slice, reflect.Array:
		compareSlice(t, v, v2)
	case reflect.Map:
		compareMap(t, v, v2)
	case reflect.Chan:
		compareChan(t, v, v2)
	default:
		comparePOD(t, v, v2)
	}
}

func comparePOD(t *testing.T, v, v2 interface{}) {
	//kind of weird. convert v2 to the same type of v1
	if v != reflect.ValueOf(v2).Convert(reflect.TypeOf(v)).Interface() {
		t.Error("not equals values: ", v, "!=", v2)
	}
}

func compareContainer(t *testing.T, v, v2 interface{}) (reflect.Value, reflect.Value) {
	val := reflect.ValueOf(v)
	val2 := reflect.ValueOf(v2)

	if val.Type() != val2.Type() {
		t.Error("Container type are not identical:", val.Type(), "!=", val2.Type())
	}
	if val.Len() != val2.Len() {
		t.Error("Container are not of the same size:", val.Len(), "!=", val2.Len())
	}
	return val, val2
}

func compareSlice(t *testing.T, v, v2 interface{}) {
	val, val2 := compareContainer(t, v, v2)
	for i := 0; i < val.Len(); i++ {
		compareValue(t, val.Index(i).Interface(), val2.Index(i).Interface())
	}
	fmt.Printf("Same slice ;)\n")
}

func compareMap(t *testing.T, v, v2 interface{}) {
	val, val2 := compareContainer(t, v, v2)
	keys := val.MapKeys()
	for _, v := range keys {
		compareValue(t, val.MapIndex(v).Interface(), val2.MapIndex(v).Interface())
	}
	fmt.Printf("Same Map ;)\n")
}

func compareChan(t *testing.T, v, v2 interface{}) {
	t.Error("Chan compare not implemented")
}

func checkValue(t *testing.T, v interface{}) {
	qival, err := toQiValue(v)
	if err != nil {
		t.Error(err)
	}
	v2, err := toGoValue(qival)
	if err != nil {
		t.Error(err)
	}
	compareValue(t, v, v2)
}

func TestPODValue(t *testing.T) {
	fmt.Println("Hello World!")
	checkValue(t, int(42))
	checkValue(t, int(-42))
	checkValue(t, uint(42))
	checkValue(t, int8(41))
	checkValue(t, int8(-41))
	checkValue(t, uint8(41))
	checkValue(t, int16(43))
	checkValue(t, int16(-43))
	checkValue(t, uint16(43))
	checkValue(t, int32(44))
	checkValue(t, int32(-44))
	checkValue(t, uint32(44))
	checkValue(t, int64(45))
	checkValue(t, int64(-45))
	checkValue(t, uint64(45))

	checkValue(t, float32(-4.6))
	checkValue(t, float32(4.6))
	checkValue(t, float64(-4.7))
	checkValue(t, float64(4.7))

	checkValue(t, string("life"))
	checkValue(t, string("42"))
}

func TestListValue(t *testing.T) {
	l := make([]int, 2)
	l[0] = 4
	l[1] = 2
	checkValue(t, l)
}

func TestMapValue(t *testing.T) {
	ma := make(map[string]int)
	ma["who"] = 42
	ma["man"] = 10
	checkValue(t, ma)

	//ch := make(chan []int)
	//checkValue(t, ch)
}
