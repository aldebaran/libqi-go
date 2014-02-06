/*
**
** Author(s):
**  - Cedric GESTES <gestes@aldebaran-robotics.com>
**
** Copyright (C) 2010, 2011, 2012 Aldebaran Robotics
*/

// BUG(r): This is not finished
package messaging

import ("fmt"
	"bytes")

type ServiceInfo struct {
	Name      string
	ServiceId uint32
	MachineId string
	ProcessId uint32
	Endpoints []string
}


func (s *ServiceInfo) QiDecode(dec *Decoder) error {
	dec.Decode(s.Name)
	dec.Decode(s.ServiceId)
	dec.Decode(s.MachineId)
	dec.Decode(s.ProcessId)
	dec.Decode(s.Endpoints)
	//todo: handle error
	return nil
}

func (s *ServiceInfo) QiEncode(enc *Encoder) error {
	enc.Encode(s.Name)
	enc.Encode(s.ServiceId)
	enc.Encode(s.MachineId)
	enc.Encode(s.ProcessId)
	enc.Encode(s.Endpoints)
	//todo: handle error
	return nil
}

//wrapped structure (to have a nice go object)
type Session struct {
	sessionSocket TransportSocket
	sessionServer TransportServer
}

// Client
func NewSession(name string) *Session {
	fmt.Printf("NewSession\n")
	//session.session = C.qi_session_create(C.CString(name))
	return &Session{}
}

func (s *Session) Connect(address string) error {
	fmt.Printf("SessionConnect: %s\n", address)
	return s.sessionSocket.Connect(address)
}

func (s *Session) DeleteSession() {
	//C.qi_client_destroy(s.psession)
}

//Run the session
func (s *Session) Run() {
	//
}

func (s *Session) Service(name string) (*Object, error) {
	msg := NewTransportMessage()
	b   := &bytes.Buffer{}
	enc := NewEncoder(b)

	fmt.Printf("Try to get Service: %s\n", name)
	enc.Encode(name)

	msg.Type     = Type_Call
	msg.Service  = Service_ServiceDirectory
	msg.Path     = Path_Main
	msg.Function = ServiceDirectoryFunction_Service
	msg.Data     = b.Bytes()

	rawData, _ := msg.QiEncode()

	fmt.Printf("Sending request...\n")

	err := s.sessionSocket.Send(rawData)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Request sent\n")

	s.sessionSocket.Receive()


	var obj Object
	//obj.pobject = C.qi_session_service(C.CString(name))
	return &obj, nil
}

func (s *Session) RegisterService(serviceName string, object interface{}) int32 {
	fmt.Printf("ServerConnect\n")
	return 1
	//C.qi_server_connect(s.pserver, C.CString(address))
}
