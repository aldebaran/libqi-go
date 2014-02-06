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

#include <qic/session.h>
#include <qic/value.h>
*/
import "C"

import (
	"errors"
	"reflect"
	"runtime"
)

// A Session manage the connection with naoqi.
// it give access to all naoqi services.
type Session struct {
	p        *C.qi_session_t
	services map[uint32]*goObjectHolder
}

func destroySession(c *Session) {
	C.qi_session_destroy(c.p)
}

// NewSession returns a new Session instance
func NewSession() *Session {
	session := &Session{C.qi_session_create(), make(map[uint32]*goObjectHolder)}
	runtime.SetFinalizer(session, destroySession)
	return session
}

func (c *Session) Dial(protocol string, address string) error {
	_, err := await(C.qi_session_connect(c.p, C.CString(protocol+"://"+address)))
	return err
}

func (c *Session) Close() error {
	_, err := await(C.qi_session_close(c.p))
	return err
}

func (s *Session) Service(name string) (AnyObject, error) {
	val, err := await(C.qi_session_get_service(s.p, C.CString(name)))
	if err != nil {
		return nil, err
	}
	v, ok := val.(AnyObject)
	if ok != true {
		return nil, errors.New("Cant convert Value to Object")
	}
	return v, nil
}

func (s *Session) RegisterService(name string, service interface{}) (uint32, error) {
	obj, err := newValueObject(reflect.ValueOf(service))
	if err != nil {
		return 0, err
	}
	if obj == nil {
		return 0, errors.New("Object is nil")
	}
	val, err := await(C.qi_session_register_service(s.p, C.CString(name), obj.cobj))
	if err != nil {
		return 0, err
	}
	s.services[val.(uint32)] = obj
	return val.(uint32), err
}

func (s *Session) UnregisterService(serviceId uint64) error {
	_, err := await(C.qi_session_unregister_service(s.p, C.uint(serviceId)))
	return err
}
