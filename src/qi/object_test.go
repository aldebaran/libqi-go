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
	"reflect"
	"testing"
)

type ServiceTest interface {
	Reply(string) string
	MapMerge(map[string]int, map[string]int) map[string]int
}

type serviceTestImpl struct{}

func (s *serviceTestImpl) Reply(val string) string {
	fmt.Println("called reply:", val)
	return val + "bim"
}

func (s *serviceTestImpl) MapMerge(a, b map[string]int) map[string]int {
	for k, v := range b {
		a[k] = v
	}
	return a
}

func TestQiObjectToGo(t *testing.T) {
	//s := NewSession()
	fmt.Println("Testing Objec")
	st := &serviceTestImpl{}
	//_, err := s.RegisterService("serviceTest", st, ServiceTest(st))
	//if err != nil {
	//	t.Fatal(err)
	//}
	serv := ServiceTest(st)
	val := reflect.ValueOf(serv)

	meth := val.Type().Method(1)
	fmt.Println("meth:", meth)

	mt := meth.Type
	fmt.Println("mt:", mt)

	fun := meth.Func
	fmt.Println("fun:", fun)

	arg1 := "titi"
	args := make([]reflect.Value, 2)
	args[0] = reflect.ValueOf(serv)
	args[1] = reflect.ValueOf(arg1)
	fun.Call(args)
	fmt.Println("done")

	//ost, err := s.Service("serviceTest")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//ost.Call("Reply", "plouf")
}
