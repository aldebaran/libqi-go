/*
**
** Author(s):
**  - Cedric GESTES <gestes@aldebaran-robotics.com>
**
** Copyright (C) 2013 Aldebaran Robotics
 */

package qi

import (
	"fmt"
	"testing"
)

var retchan = make(chan int, 10)

func onTestSig(a, b int) {
	fmt.Println("onTestSig(", a, b, ")")
	retchan <- a + b
}

func TestSubFunc(t *testing.T) {
	fmt.Println("Running test: TestSubFunc")
	s := NewSignal(func(int, int) {})
	s.Connect(onTestSig)
	s.Connect(onTestSig)
	go s.Emit(4, 2)
	r1 := <-retchan
	r2 := <-retchan
	if r1 != 6 {
		t.Error("expected", 6, "got", r1)
	}
	if r2 != 6 {
		t.Error("expected", 6, "got", r2)
	}
}

func TestSubChan(t *testing.T) {

}
