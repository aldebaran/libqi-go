/*
**
** Author(s):
**  - Cedric GESTES <gestes@aldebaran-robotics.com>
**
** Copyright (C) 2013 Aldebaran Robotics
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"qi"
	"sync"
	"time"
)

var flagEvent int

type ServiceTest interface {
	Reply(v string) string
	LastValue() string
	TestEvent() qi.Signal
}

type ServiceTest2 interface {
	Reply(v string) string
	LastValue() string
	onTestEvent(func(int, int)) qi.Signal
}

type serviceTestImpl struct {
	mutex *sync.Mutex
	value string

	bar string
	//foo chan[string]
	onEvent qi.Signal
}

func NewServiceTest() ServiceTest {
	st := &serviceTestImpl{mutex: &sync.Mutex{}, onEvent: qi.NewSignal(func(int) {})}
	return ServiceTest(st)
}

func (self *serviceTestImpl) Reply(v string) string {
	fmt.Println("Reply called:", v)
	self.mutex.Lock()
	self.value = v
	self.mutex.Unlock()
	return v + "bim"
}

func (self *serviceTestImpl) LastValue() string {
	fmt.Println("LastValue called")
	self.mutex.Lock()
	v := self.value
	self.mutex.Unlock()
	return v
}

func (self *serviceTestImpl) TestEvent() qi.Signal {
	return self.onEvent
}

func titilleMoi(st ServiceTest) {
	sti := st.(*serviceTestImpl)

	ctick := time.Tick(1 * time.Second)
	for now := range ctick {
		fmt.Println("titillage:", now)
		sti.onEvent.Emit("42")
	}
}

func main() {
	session := qi.NewSession()

	fmt.Println("## Connection")
	err := session.Dial("tcp", "127.0.0.1:9559")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(" *Connected")

	st := NewServiceTest()

	fmt.Println("## Register service")
	_, err = session.RegisterService("serviceTest", st)

	if err != nil {
		log.Fatal(err)
	}

	go titilleMoi(st)
	c := make(chan struct{})
	<-c
}

func init() {
	flag.IntVar(&flagEvent, "event", 0, "enable event")
	flag.Parse()
}
