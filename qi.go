/*
**
** Author(s):
**  - Cedric GESTES <gestes@aldebaran-robotics.com>
**
** Copyright (C) 2010, 2011, 2012 Aldebaran Robotics
*/

package qi

/*
 #include <qimessaging/c/session_c.h>
 #include <qimessaging/c/message_c.h>
*/
import "C"

import ("fmt"
 	"unsafe")

//wrapped structure (to have a nice go object)
type Session struct { psession unsafe.Pointer }
type Object  struct { pobject  unsafe.Pointer }
type Message struct { pmessage unsafe.Pointer }

// Client
func SessionCreate(name string) *Client {
	var client Client

	fmt.Printf("ClientCreate\n")
	client.pclient = C.qi_session_create(C.CString(name))
	return &client
}

func (c *Session) Connect(address string) {
	fmt.Printf("ClientConnect:", address)
	C.qi_session_connect(c.pclient, C.CString(address))
}

func (c *Session) Destroy() {
	C.qi_client_destroy(c.pclient)
}

func (c *Object) Call(methodSignature string, params *Message, ret *Message) {
	fmt.Printf("ClientCall\n")
	C.qi_session_call(c.pclient, C.CString(methodSignature), params.pmessage, ret.pmessage)
}

func (c *Object) Go(methodSignature string, params *Message, ret *Message) {
	fmt.Printf("ClientCall\n")
	C.qi_session_call(c.pclient, C.CString(methodSignature), params.pmessage, ret.pmessage)
}


// Message
func MessageCreate() *Message {
	var message Message

	fmt.Printf("MessageCreate\n")
	message.pmessage = C.qi_message_create()
	return &message
}

func (m *Message) WriteChar(c byte) {
	C.qi_message_write_char(m.pmessage, C.char(c))
}

func (m *Message) ReadChar() byte {
	return (byte)(C.qi_message_read_char(m.pmessage))
}

func (m *Message) WriteInt(i int32) {
	C.qi_message_write_int(m.pmessage, C.int(i))
}

func (m *Message) ReadInt() int32 {
	return (int32)(C.qi_message_read_int(m.pmessage))
}

func (m *Message) WriteFloat(i float32) {
	C.qi_message_write_float(m.pmessage, C.float(i))
}

func (m *Message) ReadFloat() float32 {
	return (float32)(C.qi_message_read_float(m.pmessage))
}

func (m *Message) WriteDouble(i float64) {
	C.qi_message_write_double(m.pmessage, C.double(i))
}

func (m *Message) ReadDouble() float64 {
	return (float64)(C.qi_message_read_double(m.pmessage))
}

func (m *Message) WriteString(s string) {
	C.qi_message_write_string(m.pmessage, C.CString(s))
}

func (m *Message) ReadString() string {
	return C.GoString(C.qi_message_read_string(m.pmessage))
}


// Server
func ServerCreate(name string) *Server {
	var server Server

	fmt.Printf("ServerCreate\n")
	server.pserver = C.qi_server_create(C.CString(name))
	return &server
}

func (s *Server) Connect(address string) {
	fmt.Printf("ServerConnect\n")
	C.qi_server_connect(s.pserver, C.CString(address))
}

func (s *Server) Destroy() {
	C.qi_server_destroy(s.pserver)
}
