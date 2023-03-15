package tcprpc

import (
	"net"

	"fsky.pro/fsrpcs"
	//	"fsky.pro/fsrpcs/tcprpc"
)

// -------------------------------------------------------------------
// server
// -------------------------------------------------------------------
type S_Service1 struct {
	fsrpcs.I_Service
}

func (*S_Service1) Name() string {
	return "service1"
}

func (*S_Service1) Hello(req string, reply *string) error {
	*reply = "hello my name is server1"
	return nil
}

// ---------------------------------------------------------
type S_Service2 struct {
	fsrpcs.I_Service
}

func (*S_Service2) Name() string {
	return "service1"
}

func (*S_Service2) Hello(req string, reply *string) error {
	*reply = "hello my name is server2"
	return nil
}

// ---------------------------------------------------------
// 一个 server 处理多个 service
func sever_main() {
	server := NewServer()
	server.Register(new(S_Service1))
	server.Register(new(S_Service2))
	addr, _ := net.ResolveTCPAddr("tcp", ":8080")
	if _, err := server.Listen("tcp", addr); err == nil {
		panic("listen error!")
	}
	server.Serve(nil)
}

// -------------------------------------------------------------------
// client
// -------------------------------------------------------------------
// 对应服务器端的 Service1
type S_Client1 struct {
	*fsrpcs.S_Client
}

func NewClient1(proxy fsrpcs.I_ServiceProxy) *S_Client1 {
	return &S_Client1{
		S_Client: fsrpcs.NewClient("service1", proxy),
	}
}

func (this *S_Client1) Hello(req string, reply *string) error {
	return this.S_Client.Call("Hello", req, reply)
}

// ---------------------------------------------------------
// 对应服务器端的 Service2
type S_Client2 struct {
	*fsrpcs.S_Client
}

func NewClient2(proxy fsrpcs.I_ServiceProxy) *S_Client2 {
	return &S_Client2{
		S_Client: fsrpcs.NewClient("service2", proxy),
	}
}

func (this *S_Client2) Hello(req string, reply *string) error {
	return this.S_Client.Call("Hello", req, reply)
}

// ---------------------------------------------------------
// 一个 proxy 对应多个 client
// 一个 client 对应一个 service
func client_main() {
	addr, _ := net.ResolveTCPAddr("tcp", "localhost:8080")
	proxy := NewServerProxy("tcp", addr)
	client1 := NewClient1(proxy)
	client2 := NewClient2(proxy)

	if _, err := proxy.Dial(); err != nil {
		panic("dial server faili!")
	}

	var reply1, reply2 string
	client1.Hello("client1", &reply1)
	client2.Hello("client2", &reply2)
}
