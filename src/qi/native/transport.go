
package messaging

import ("strings"
	"net"
	"fmt"
	"log"
	"sync/atomic"
	"bytes")

const (
	Type_None  uint32 = 0
	Type_Call  uint32 = 1
	Type_Reply uint32 = 2
	Type_Event uint32 = 3
	Type_Error uint32 = 4
	)

const (
	Service_Server           uint32 = 0
	Service_ServiceDirectory uint32 = 1
)

const (
	Path_None                uint32 = 0
	Path_Main                uint32 = 1
)

const (
	Function_MetaObject      uint32 = 0
)

const (
	ServiceDirectoryFunction_Service           uint32 = 1
	ServiceDirectoryFunction_Services	   uint32 = 2
	ServiceDirectoryFunction_RegisterService   uint32 = 3
	ServiceDirectoryFunction_UnregisterService uint32 = 4
	ServiceDirectoryFunction_ServiceReady	   uint32 = 5
)

type TransportMessage struct {
	id          uint32
	Size        uint32
	Type        uint32
	Service     uint32
	Path        uint32
	Function    uint32
	Data      []byte
}

//uniq message id
var messageId uint32 = 0

func NewTransportMessage() *TransportMessage {
	msg    := &TransportMessage{}
	msg.id  = atomic.AddUint32(&messageId, 1)
	return msg
}

func NewTransportMessageByte(raw []byte) *TransportMessage {
	msg := &TransportMessage{}

	b   := bytes.NewBuffer(raw)
	dec := NewDecoder(b)

	var magic    uint32
	var reserved uint32

	dec.Decode(&magic)
	dec.Decode(&msg.id)
	dec.Decode(&msg.Size)
	dec.Decode(&msg.Type)
	dec.Decode(&msg.Service)
	dec.Decode(&msg.Path)
	dec.Decode(&msg.Function)
	dec.Decode(&reserved)
	msg.Data = b.Bytes()
	// if int(msg.Size) != len(msg.Data) {
	// 	log.Fatal("bad message size %d != %d\n", msg.Size, len(msg.Data))
	// }
	return msg
}

func (t *TransportMessage) QiEncode() ([]byte, error) {
	var b     bytes.Buffer
	enc := NewEncoder(&b)
	var magic uint32 = 0x42adde42

	enc.Encode(uint32(magic))
	enc.Encode(uint32(t.id))
	enc.Encode(uint32(len(t.Data)))
	enc.Encode(uint32(t.Type))
	enc.Encode(uint32(t.Service))
	enc.Encode(uint32(t.Path))
	enc.Encode(uint32(t.Function))
	enc.Encode(uint32(0))
	b.Write(t.Data)
	return b.Bytes(), nil
}

func (t *TransportMessage) QiDecode(b []byte) error {
	return nil
}

//////////////////////////////// TransportSocket

type TransportSocket struct {
	address  string;
	conn     net.Conn;
}

func urlToProtoAddress(url string) (string, string, error) {
	res := strings.SplitN(url, "://", 2)
	if len(res) < 2 {
		return "", "", fmt.Errorf("split failed %d\n", len(res))
	}
	return res[0], res[1], nil
}

func NewTransportSocket(conn net.Conn) *TransportSocket {
	return &TransportSocket{conn: conn}
}

func (tso *TransportSocket) Connect(address string) error {
	proto, addr, err := urlToProtoAddress(address)
	if err != nil {
		return err
	}

	fmt.Printf("Protocol: %s\n", proto)
	fmt.Printf("Address: %s\n", addr)

	conn, err := net.Dial(proto, addr)
	if err != nil {
		return err
		// handle error
	}
	tso.address = address;
	tso.conn = conn
	fmt.Printf("Connected to %s\n", address)
	return nil
	//fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	//status, err := bufio.NewReader(conn).ReadString('\n')
// ...
}

func (tso *TransportSocket) Send(data[] byte) error {
	sz, err := tso.conn.Write(data)
	fmt.Printf("Sent: %d\n", sz)
	return err
	//fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	//status, err := bufio.NewReader(conn).ReadString('\n')
}

func (tso *TransportSocket) Receive() *TransportMessage {
	b := make([]byte, 32)
	fmt.Printf("Reading..\n")
	sz, err := tso.conn.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Read size: %d\n", sz)
	msg := NewTransportMessageByte(b)
	fmt.Printf("Msg size: %d\n", msg.Size)
	b2 := make([]byte, msg.Size)
	sz, err = tso.conn.Read(b2)
	msg.Data = b2
	return msg
}


type TransportServer struct {
	listenAddress string;
}

func NewTransportServer() *TransportServer {
	var ts TransportServer
	return &ts
}

func (tso *TransportServer) Listen(address string, rq chan *TransportSocket) error {
	proto, addr, err := urlToProtoAddress(address)
	if err != nil {
		return fmt.Errorf("Bad url %s", address)
	}

	fmt.Printf("Protocol: %s\n", proto)
	fmt.Printf("Address: %s\n", addr)

	// Listen on TCP port 2000 on all interfaces.
	l, err := net.Listen(proto, addr)
	if err != nil {
		log.Fatal(err)
		return err
	}
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		fmt.Printf("Accepted a new connection\n")
		if err != nil {
			log.Fatal(err)
		}
		rq <- NewTransportSocket(conn)
	}
	return nil
}
