/*
**
** Author(s):
**  - Cedric GESTES <gestes@aldebaran-robotics.com>
**
** Copyright (C) 2010, 2011, 2012 Aldebaran Robotics
*/

// BUG(r): This is not finished
package messaging

/*
 #include <qimessaging/c/object_c.h>
*/
//import "C"

import ("fmt"
 	"unsafe")

//{IssI}{IsI}I
type MetaMethod struct {
	Signature string
	SigReturn string
	Uid       uint32
}

type MetaEvent struct {
	Signature string
	Uid       uint32
}

type MetaObject struct {
	MetaMethods map[uint32]MetaMethod
	MetaEvents  map[uint32]MetaEvent
	LastUid       uint32
}



type Object struct { object unsafe.Pointer }

/** Object
 *
 * Call
 * Go
 * Subscribe
 * Publish
 */
type Call struct {
    MethodSignature string      // The name of the service and method to call.
    Args            interface{} // The argument to the function (*struct).
    Reply           interface{} // The reply from the function (*struct).
    Error           error       // After completion, the error status.
    Done            chan *Call  // Strobes when call is complete.
}
func (c *Object) Call(methodSignature string, args interface{}, reply interface {}) error {
	fmt.Printf("ClientCall\n")
	//C.qi_session_call(c.pobject, C.CString(methodSignature), params.pmessage, ret.pmessage)
	return nil
}

func (c *Object) Go(methodSignature string, args interface{}, reply interface{}, done chan *Call) *Call {
	fmt.Printf("Object::Go")
	//C.qi_session_call(c.pobject, C.CString(methodSignature), params.pmessage, ret.pmessage)
	var call Call
	return &call
}

func (c *Object) Subscribe(signalSignature string, events chan interface{}) {
	//create a future on the signal
	//create a go routine that wait in loop for the future
	//go handle channel
}

func (c *Object) Publish(slotSignature string, events chan interface{}) {
	//create goroutine that wait on the channel
	//and send to the slot
}
