
package messaging

import "testing"
import "fmt"


type TestIS struct {
	A int
	B string
}

func TestTakeSignatre(t *testing.T) {
	var i8 int8 = 2
	var i16 int16 = 2
	var i32 int32 = 2
	var i64 int64 = 2

	var ui8	 uint8 = 2
	var ui16 uint16 = 2
	var ui32 uint32 = 2
	var ui64 uint64 = 2

	var f32 float32 = 3.4
	var f64 float64 = 3.4

	msi := make(map[string]int)
	tis := TestIS{23, "skidoo"}

	fmt.Println("Signature:", TakeSignature(i8))
	fmt.Println("Signature:", TakeSignature(i16))
	fmt.Println("Signature:", TakeSignature(i32))
	fmt.Println("Signature:", TakeSignature(i64))

	fmt.Println("Signature:", TakeSignature(ui8))
	fmt.Println("Signature:", TakeSignature(ui16))
	fmt.Println("Signature:", TakeSignature(ui32))
	fmt.Println("Signature:", TakeSignature(ui64))

	fmt.Println("Signature:", TakeSignature(f32))
	fmt.Println("Signature:", TakeSignature(f64))
	fmt.Println("Signature:", TakeSignature(msi))
	fmt.Println("Signature:", TakeSignature(tis))
	//t.Errorf("Sqrt(%v) = %v, want %v", in, x, out)
}
