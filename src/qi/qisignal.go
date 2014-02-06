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
#include <qic/future.h>
#include <stdlib.h>
qi_future_t* go_object_signal_connect(qi_object_t *obj, char *signame, void *userdata);
*/
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

type qiProxySignal struct {
	gosig Signal //the go signal
	cobj  *C.qi_object_t
	name  *C.char
	cid   uint64
}

//export go_object_signal_callback
func go_object_signal_callback(cval unsafe.Pointer, qps unsafe.Pointer) {
	//send the value received to the channel
	proxy := ((*qiProxySignal)(qps))
	val := (*C.qi_value_t)(cval)

	gv, err := toGoValue(val)
	if err != nil {
		fmt.Println("Discarding signal value because it's invalid")
	}
	proxy.gosig.Emit(gv.([]interface{})...)
	//TODO: destroy val?
}

func newProxySignal(cobj *C.qi_object_t, name string) Signal {
	//TODO: handle types here
	sig := &qiProxySignal{gosig: NewSignal(func(...interface{}) {}), cobj: cobj, name: C.CString(name)}
	id, err := await(C.go_object_signal_connect(cobj, sig.name, unsafe.Pointer(sig)))
	if err != nil {
		panic(err)
	}
	sig.cid = uint64(id.(uint64))
	runtime.SetFinalizer(sig, destroyProxySignal)
	return Signal(sig)
}

//Here we assume p.cobj is alive
func destroyProxySignal(p *qiProxySignal) {
	C.free(unsafe.Pointer(p.name))
	C.qi_object_signal_disconnect(p.cobj, C.ulonglong(p.cid))
}

//Only emit to the C signal. the C->GO callback will call the gosig as appropriated
func (s *qiProxySignal) Emit(v ...interface{}) error {
	params, err := createTupleValue(v...)
	if err != nil {
		return err
	}
	ok := C.qi_object_post(s.cobj, s.name, params)
	if ok != 1 {
		return errors.New("Object.Post returned an error")
	}
	return nil
}

func (s *qiProxySignal) Connect(chanOrFun interface{}) SignalSubscriberId {
	return s.gosig.Connect(chanOrFun)
}

func (s *qiProxySignal) Disconnect(id SignalSubscriberId) {
	s.gosig.Disconnect(id)
}
