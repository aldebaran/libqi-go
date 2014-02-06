

package messaging

import "testing"
import "fmt"
import "bytes"


func TestSimpleTypes(t *testing.T) {
	var b    bool  = true
	var i8   int8  = 1
	var i16  int16 = 2
	var i32  int32 = 3
	var i64  int64 = 4

	var ui8	 uint8  = 5
	var ui16 uint16 = 6
	var ui32 uint32 = 7
	var ui64 uint64 = 8

	var f32 float32 = 9.9
	var f64 float64 = 10.10

	var rb    bool
	var ri8	  int8
	var ri16  int16
	var ri32  int32
	var ri64  int64

	var rui8  uint8
	var rui16 uint16
	var rui32 uint32
	var rui64 uint64

	var rf32 float32
	var rf64 float64


	var err error
	var buf bytes.Buffer
	qie := NewEncoder(&buf)


	if err = qie.Encode(b	); err != nil { t.Error(err) }
	if err = qie.Encode(i8	); err != nil { t.Error(err) }
	if err = qie.Encode(i16 ); err != nil { t.Error(err) }
	if err = qie.Encode(i32 ); err != nil { t.Error(err) }
	if err = qie.Encode(i64 ); err != nil { t.Error(err) }
	if err = qie.Encode(ui8 ); err != nil { t.Error(err) }
	if err = qie.Encode(ui16); err != nil { t.Error(err) }
	if err = qie.Encode(ui32); err != nil { t.Error(err) }
	if err = qie.Encode(ui64); err != nil { t.Error(err) }
	if err = qie.Encode(f32 ); err != nil { t.Error(err) }
	if err = qie.Encode(f64 ); err != nil { t.Error(err) }

	qid := NewDecoder(&buf)
	if err = qid.Decode(&rb	  ); err != nil { t.Error(err) }
	if err = qid.Decode(&ri8  ); err != nil { t.Error(err) }
	if err = qid.Decode(&ri16 ); err != nil { t.Error(err) }
	if err = qid.Decode(&ri32 ); err != nil { t.Error(err) }
	if err = qid.Decode(&ri64 ); err != nil { t.Error(err) }
	if err = qid.Decode(&rui8 ); err != nil { t.Error(err) }
	if err = qid.Decode(&rui16); err != nil { t.Error(err) }
	if err = qid.Decode(&rui32); err != nil { t.Error(err) }
	if err = qid.Decode(&rui64); err != nil { t.Error(err) }
	if err = qie.Encode(&rf32) ; err != nil { t.Error(err) }
	if err = qie.Encode(&rf64) ; err != nil { t.Error(err) }

	if b    != rb	 { t.Errorf("Value decoded mismatch") }
	if i8	!= ri8	 { t.Errorf("Value decoded mismatch") }
	if i16	!= ri16	 { t.Errorf("Value decoded mismatch") }
	if i32	!= ri32	 { t.Errorf("Value decoded mismatch") }
	if i64	!= ri64	 { t.Errorf("Value decoded mismatch") }
	if ui8	!= rui8	 { t.Errorf("Value decoded mismatch") }
	if ui16 != rui16 { t.Errorf("Value decoded mismatch") }
	if ui32 != rui32 { t.Errorf("Value decoded mismatch") }
	if ui64 != rui64 { t.Errorf("Value decoded mismatch") }

}

func TestStrings(t *testing.T) {
	var s1, s2 string

	s1 = "mes couilles en string"
	var err error
	var buf bytes.Buffer
	qie := NewEncoder(&buf)
	if err = qie.Encode(s1	); err != nil { t.Error(err) }

	qid := NewDecoder(&buf)
	if err = qid.Decode(&s2	  ); err != nil { t.Error(err) }

	if s1 != s2 { t.Error("Value decoded mismatch") }

}


type TestStruct struct {
	A int
	B string
}

func TestSTLTypes(t *testing.T) {
	//yes the name of the test if funny
	//t.Errorf("Sqrt(%v) = %v, want %v", in, x, out)
	fmt.Printf("claclk")
}
