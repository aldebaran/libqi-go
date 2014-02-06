/*
**
** Author(s):
**  - Cedric GESTES <gestes@aldebaran-robotics.com>
**
** Copyright (C) 2013 Aldebaran Robotics
 */

package qi

/* wrap Go object into a qiObject
 */

/*
#cgo CFLAGS: -I/home/ctaf/src/qi-master/lib/libqi -I/home/ctaf/src/qi-master/lib/libqimessaging/c/
#cgo LDFLAGS: -L/home/ctaf/src/qi-master/lib/libqimessaging/build-sys-linux-x86_64/sdk/lib/ -lqic

#include <qic/value.h>
#include <qic/object.h>

unsigned long long go_object_builder_advertise_method(qi_object_builder_t *ob, char *signame, void *callback);
*/
import "C"

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

//used to hold go structure given to C alive.
type goObjectHolder struct {
	meths []*methodHolder
	sigs  []*signalHolder
	obj   reflect.Value
	cobj  *C.qi_object_t //c wrapped object
}

type methodHolder struct {
	fun reflect.Value //functor
	obj reflect.Value //parent object
}

type signalHolder struct {
	cobj **C.qi_object_t //pointer on pointer because we dont know cobj when creating the struct
	name string
	sig  string
	id   uint64
}

func qiMethodType(meth reflect.Method) string {
	if meth.Type.NumOut() != 1 {
		return "func"
	}
	ret := meth.Type.Out(0)
	if ret.Kind() != reflect.Interface {
		return "func"
	}
	//todo: someday it will have a better name
	if ret.PkgPath() != "qi" {
		return "func"
	}
	if ret.Name() == "Signal" {
		return "signal"
	} else if ret.Name() == "Property" {
		return "property"
	}
	fmt.Println("func iface type name:", ret.Name())
	return "func"
}

func convertRetValue(retval []reflect.Value, toreturn *C.qi_value_t) error {
	switch len(retval) {
	case 0:
		return nil
	case 1:
		if _, err := valueSet(toreturn, retval[0]); err != nil {
			return err
		}
	default:
		for k, rv := range retval {
			tr := C.qi_value_tuple_get(toreturn, C.uint(k))
			if tr == nil {
				return errors.New("Cant convert retval")
			}
			if _, err := valueSet(tr, rv); err != nil {
				return err
			}

		}
	}
	return nil
}

//export go_object_call_callback
func go_object_call_callback(sig, value, toreturn, userdata unsafe.Pointer) {
	fun := *(*methodHolder)(userdata)
	fmt.Println("Should call method here:", fun)
	vals, err := toGoValue((*C.qi_value_t)(value))
	if err != nil {
		fmt.Println("cant convert args..")
		//TODO set ret as appropriated
		return
	}
	rvals := vals.([]interface{})
	ln := len(rvals)

	args := make([]reflect.Value, ln+1)
	args[0] = fun.obj
	for k, v := range rvals {
		fmt.Println("val:", v)
		args[k+1] = reflect.ValueOf(v)
	}
	fmt.Println("calllllinnnnnnng:", fun, args)
	ret := fun.fun.Call(args)
	fmt.Println("call done")
	if convertRetValue(ret, (*C.qi_value_t)(toreturn)) != nil {
		fmt.Println("Cant convert val:", err)
		//TODO set the error
	}
}

func advertiseMethod(oh *goObjectHolder, ob *C.qi_object_builder_t, mh *methodHolder, meth reflect.Method) error {
	sig, err := methodSignature(meth)
	if err != nil {
		return err
	}
	fmt.Println("Found method: ", sig)
	mh.fun = meth.Func
	mh.obj = oh.obj
	ok := C.go_object_builder_advertise_method(ob, C.CString(sig), unsafe.Pointer(mh))
	if ok == 0 {
		return errors.New("Cant advertise method")
	}
	return nil
}

func signalForwardToQiObject(sh *signalHolder, v ...interface{}) {
	params, err := createTupleValue(v...)
	if err != nil {
		fmt.Println("error in forward sig:", err)
		return
	}
	fmt.Println("calling signal:", sh.name+"::"+sh.sig)
	fmt.Println("cobj:", sh.cobj)
	if sh.cobj == nil || *sh.cobj == nil {
		fmt.Println("forwardSig: cobj is nil")
		return
	}
	fmt.Println("send data of sig:", C.GoString(C.qi_value_get_signature(params, 0)))
	C.qi_object_post(*sh.cobj, C.CString(sh.name), params)
}

func advertiseSignal(oh *goObjectHolder, ob *C.qi_object_builder_t, sh *signalHolder, meth reflect.Method) error {
	//TODO: compute the real signature of the signal
	ok := C.qi_object_builder_advertise_signal(ob, C.CString(meth.Name), C.CString("m"))
	if ok == 0 {
		return errors.New("Cant advertise method")
	}
	sh.cobj = &oh.cobj
	sh.name = meth.Name
	sh.sig = "m"
	sh.id = uint64(ok)
	outs := meth.Func.Call([]reflect.Value{oh.obj})
	sig := outs[0].Interface().(Signal)
	sig.Connect(func(v ...interface{}) { signalForwardToQiObject(sh, v...) })
	return nil
}

func newValueObject(val reflect.Value) (*goObjectHolder, error) {
	l := val.NumMethod()
	//TODO: allocate only what is needed...
	oh := &goObjectHolder{
		obj:   val,
		meths: make([]*methodHolder, l),
		sigs:  make([]*signalHolder, l),
	}

	ob := C.qi_object_builder_create()
	for i := 0; i < l; i++ {
		meth := val.Type().Method(i)
		qimethtype := qiMethodType(meth)
		switch qimethtype {
		case "func":
			oh.meths[i] = &methodHolder{}
			if err := advertiseMethod(oh, ob, oh.meths[i], meth); err != nil {
				return nil, err
			}
		case "signal":
			oh.sigs[i] = &signalHolder{}
			if err := advertiseSignal(oh, ob, oh.sigs[i], meth); err != nil {
				return nil, err
			}
		case "property":
			fmt.Printf("Not Implemented: properties")
			//	C.go_object_builder_advertise_property(C.CString(meth.Name()))
		}
	}
	oh.cobj = C.qi_object_builder_get_object(ob)
	if oh.cobj == nil {
		return nil, errors.New("object builder returned nil")
	}
	return oh, nil
}
