
package main

import "fmt"
import qi "../_obj/qi"

func main() {
	fmt.Printf("Creating a connection\n")
	client := qi.ClientCreate("goclient")
	client.Connect("127.0.0.1:5555")

	msg := qi.MessageCreate()
	ret := qi.MessageCreate()

	msg.WriteString("master.listServices::{ss}:")
	client.Call("master.listServices::{ss}:", msg, ret)
	sz := ret.ReadInt()
	for i := 0; i < (int)(sz); i++ {
		k := ret.ReadString()
		v := ret.ReadString()
		fmt.Printf("%s : %s\n", k, v)
	}
}

