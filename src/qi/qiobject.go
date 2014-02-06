/*
**
** Author(s):
**  - Cedric GESTES <gestes@aldebaran-robotics.com>
**
** Copyright (C) 2013 Aldebaran Robotics
 */

package qi

/*
#cgo CFLAGS: -I/home/ctaf/src/qi-master/lib/libqi -I/home/ctaf/src/qi-master/lib/libqimessaging/c/
#cgo LDFLAGS: -L/home/ctaf/src/qi-master/lib/libqimessaging/build-sys-linux-x86_64/sdk/lib/ -lqic

#include <qic/value.h>
#include <qic/object.h>

unsigned long long go_object_signal_connect(qi_object_t *obj, char *signame, void *userdata);
*/
import "C"

import (
	"errors"
	"fmt"
	"reflect"
)

//Public Signal Subscriber
type SignalSubscriber struct {
	C      <-chan []interface{}
	Quit   chan<- int
	Id     int64
	writer chan<- []interface{}
}

//Public Interface for all Objects
type AnyObject interface {
	Call(name string, params ...interface{}) (interface{}, error)
	Post(name string, params ...interface{}) error

	Signal(name string) (Signal, error)
	Property(name string) (Property, error)
}

/* Proxy object. Give qiObject to the Go word
 */

type qiObject struct {
	p    *C.qi_object_t
	sigs map[string]Signal
}

func toObject(v *C.qi_value_t) (AnyObject, error) {
	if v == nil {
		return nil, errors.New("Cant convert Value to Object: value is nil")
	}
	ret := C.qi_value_object_get(v)
	if ret == nil {
		return nil, errors.New("Cant convert Value to Object")
	}
	//TODO: setFinalizer
	return AnyObject(&qiObject{ret, make(map[string]Signal)}), nil
}

func (o *qiObject) Call(methodname string, v ...interface{}) (interface{}, error) {
	fmt.Printf("Object.Call\n")
	vals, err := createTupleValue(v...)
	if err != nil {
		return nil, err
	}
	val, err := await(C.qi_object_call(o.p, C.CString(methodname), vals))
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (o *qiObject) Post(methodname string, v ...interface{}) error {
	fmt.Printf("Object.Post\n")
	vals, err := createTupleValue(v...)
	if err != nil {
		return err
	}
	ret := C.qi_object_post(o.p, C.CString(methodname), vals)
	if ret == 0 {
		return nil
	}
	return errors.New(fmt.Sprintf("Post error: %d", ret))
}

func (o *qiObject) PrintMetaObject() {
	met, err := toGoValue(C.qi_object_get_metaobject(o.p))
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("Metaobject: ", met)
}

func (o *qiObject) Signal(name string) (Signal, error) {
	//TODO: not threadsafe
	v, ok := o.sigs[name]
	if ok {
		return v, nil
	}
	nsig := newProxySignal(o.p, name)
	o.sigs[name] = nsig
	return nsig, nil
}

func (o *qiObject) Property(name string) (Property, error) {
	fmt.Println("Not implemented: Property(name)")
	return Property(nil), nil
}

func valueSetObject(v *C.qi_value_t, val reflect.Value) (*C.qi_value_t, error) {
	return nil, errors.New("Cant set object value")
}
