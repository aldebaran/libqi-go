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

#include <qic/future.h>
#include <qic/value.h>

void go_async_waiter_callback(void* fut, void *userdata);
void go_future_add_callback(qi_future_t *fut, void *);

*/
import "C"

import (
	"errors"
	"log"
	"sync"
	"unsafe"
)

//########### Future

//Future are useless in GO, in GO asynchronous is futile.
//So use coroutine and avoid future. (just use a synchronous call instead)

//export go_async_waiter_callback
func go_async_waiter_callback(fut unsafe.Pointer, userdata unsafe.Pointer) {
	ai := *((*awaitInfo)(userdata))

	//notify await that we have been called
	ai.mutex.Lock()
	ai.cond.Signal()
	ai.mutex.Unlock()

	//we do not unregister from the future...
	//because it wont be called again.
	//we assume this is safe
}

type awaitInfo struct {
	mutex sync.Mutex
	cond  *sync.Cond
}

func await(f *C.qi_future_t) (interface{}, error) {
	//Deal with the coroutine engine.  Is that faster than not having it?
	//is it?
	//ai will be alive til we wait for the condition,
	//so it's safe to give it to the C future (we will detroy the cb then)
	ai := &awaitInfo{}
	ai.cond = sync.NewCond(&ai.mutex)
	//take the address of go_async_waiter_callback in a C function
	C.go_future_add_callback(f, unsafe.Pointer(ai))
	//Wait for the future to be notifieds
	ai.mutex.Lock()
	ai.cond.Wait()
	ai.mutex.Unlock()

	if int(C.qi_future_has_error(f, C.QI_FUTURETIMEOUT_INFINITE)) == 1 {
		log.Print("await error:", C.GoString(C.qi_future_get_error(f)))
		return nil, errors.New(C.GoString(C.qi_future_get_error(f)))
	}
	val := C.qi_future_get_value(f)
	ret, err := toGoValue(val)
	C.qi_value_destroy(val)
	C.qi_future_destroy(f)
	return ret, err
}
