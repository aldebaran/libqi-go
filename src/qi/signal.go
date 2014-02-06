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
	"sync"
)

type SignalSubscriberId uint

//Public Signal Interface
type Signal interface {
	Connect(chanOrFun interface{}) SignalSubscriberId
	Disconnect(id SignalSubscriberId)
	Emit(...interface{}) error //or maybe only interface{}
}

//internal signal structure
type signalImpl struct {
	sigType reflect.Type
	mutex   sync.Mutex
	subs    map[SignalSubscriberId]reflect.Value //eithers function or channel
	nextId  SignalSubscriberId
}

func NewSignal(signature interface{}) Signal {
	t := reflect.TypeOf(signature)
	if t.Kind() != reflect.Func {
		panic(fmt.Sprint("signature is not a function but has type: ", t))
	}
	return Signal(&signalImpl{sigType: t, subs: make(map[SignalSubscriberId]reflect.Value)})
}

func (s *signalImpl) Connect(chanOrFunc interface{}) (id SignalSubscriberId) {
	value := reflect.ValueOf(chanOrFunc)
	switch value.Type().Kind() {
	case reflect.Chan:
		//TODO: if s.sigType is of only one type: check elem
		//TODO: if s.sigType is a tuple only support chan []interface{} and chan interface{} ? (and of matching struct?)
	case reflect.Func:
		//TODO: check for assignable
		//TODO: check for ...interface{}
		//if !value.Type().AssignableTo(s.sigType) {
		//	panic(fmt.Sprint("func is not of the good type, got:", value.Type(), "expected:", s.sigType))
		//}
	default:
		panic("chanOrFunc is not a channel nor a function")
	}

	//add the subscriber to the map
	s.mutex.Lock()
	id = s.nextId
	s.nextId++
	s.subs[id] = value
	s.mutex.Unlock()
	return
}

//This is safe even when the id do not exits
func (s *signalImpl) Disconnect(id SignalSubscriberId) {
	s.mutex.Lock()
	delete(s.subs, id)
	s.mutex.Unlock()
}

func (s *signalImpl) Emit(v ...interface{}) error {
	fmt.Println("Emiting:", v)
	gv := make([]reflect.Value, len(v))
	for i := 0; i < len(v); i++ {
		gv[i] = reflect.ValueOf(v[i])
	}
	for _, f := range s.subs {
		switch f.Type().Kind() {
		case reflect.Chan:
			f.Send(reflect.ValueOf(v))
		case reflect.Func:
			go f.Call(gv)
		}
	}
	return nil
}
