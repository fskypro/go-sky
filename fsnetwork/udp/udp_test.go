package udp

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func server(wg *sync.WaitGroup) {
	info, _ := NewUDPInfo("127.0.0.1", 8881, 1024)
	svr, err := NewServer(info)
	if err != nil {
		fmt.Println("server error: ", err)
		wg.Done()
		return
	}

	svr.OnReceived = func(err error, cli *S_RemoteClient, data []byte) {
		if err != nil {
			fmt.Println("receive client error:", err)
			return
		}
		cli.Reply([]byte(fmt.Sprintf("receive %q ok!", string(data))))
		fmt.Println("receive from client: ", string(data))
		if string(data) == "hello 1" {
			svr.Close()
		}
	}
	svr.OnClosed = func() {
		fmt.Println("server has closed!")
		wg.Done()
	}
	go svr.Serve()
}

func client(wg *sync.WaitGroup) {
	defer wg.Done()
	info, _ := NewUDPInfo("127.0.0.1", 8881, 1024)
	client := NewClient(info)
	if err := client.Dial(); err != nil {
		fmt.Println("client err: ", err)
		return
	}
	go client.Serve()

	client.OnReceived = func(err error, data []byte) {
		if err != nil {
			fmt.Println("receive from server error: ", err)
			return
		}
		fmt.Println("response form server: ", string(data))
	}

	count := 10
	for count > 0 {
		time.Sleep(time.Second)
		msg := fmt.Sprintf("hello %d", count)
		client.Send([]byte(msg))
		count -= 1
	}
	client.Close()
}

func TestUDPServer(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	go server(&wg)
	go client(&wg)
	wg.Wait()
}
