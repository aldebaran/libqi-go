
package messaging

import "testing"

// func serve(ts *TransportServer) {
// 	rq := make(chan *transport.TransportSocket)
// 	go ts.Listen("tcp://127.0.0.1:5555", rq)
// 	for {
// 		//read a request
// 		ts <- rq
// 		//start a go routine to handle the new connection

// 	}
// }

func TestTransport(t *testing.T) {

	// ts = transport.NewTransportServer()
	// go serve(ts)


	// tso = transport.NewTransportSocket()
	// tso.Connect("tcp://127.0.0.1:5555")
	// tso.Send("plouf")
	//rec = tso.Receive()
	t.Errorf("No test implemented")
}
