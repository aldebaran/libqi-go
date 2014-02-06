/*
**
** Author(s):
**  - Cedric GESTES <gestes@aldebaran-robotics.com>
**
** Copyright (C) 2013 Aldebaran Robotics
 */

package main

import "fmt"
import "log"
import "qi"
import "flag"

var flagEvent int

func main() {
	session := qi.NewSession()
	fmt.Println("## Connecting...")
	err := session.Dial("tcp", "127.0.0.1:9559")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(" * Connected")

	fmt.Println("## Getting service")
	srv, err := session.Service("serviceTest")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(" *Got service")

	if flagEvent == 0 {
		fmt.Println("## Calling method")
		ret, err := srv.Call("reply", "plat")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(" * Recv: %s\n", ret)
	}

	if flagEvent == 1 {
		fmt.Println("## Subscribing to event")
		ch := make(chan []interface{})
		sig, err := srv.Signal("testEvent")
		if err != nil {
			log.Fatal(err)
		}
		subid := sig.Connect(ch)
		for ev := range ch {
			fmt.Printf("Recv: %s\n", ev[0].(string))
		}
		sig.Disconnect(subid)
	}
}

func init() {
	flag.IntVar(&flagEvent, "event", 0, "enable event")
	flag.Parse()
}
