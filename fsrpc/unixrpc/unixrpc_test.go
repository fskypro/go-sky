package unixrpc

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/rpc/jsonrpc"
	"sync"
	"testing"
	"time"
)

// -------------------------------------------------------------------
// proto
// -------------------------------------------------------------------
// 请求参数
type S_Arg struct {
	Order   int
	Message string `json:"msg"`
}

// 请求返回
type S_Ret struct {
	OK bool `json:"ok"`
}

// 交互代理
type S_Register struct {
	*S_CallWrap[S_Register, S_Arg, S_Ret]
}

func newRegister() *S_Register {
	return &S_Register{
		S_CallWrap: NewCallWrap[S_Register, S_Arg, S_Ret](),
	}
}

// -------------------------------------------------------------------
// server
// -------------------------------------------------------------------
type S_Receiver struct{}

func (this *S_Receiver) OnCalled(arg *S_Arg, ret *S_Ret) error {
	*ret = S_Ret{true}
	n := rand.Intn(10) + 4
	wait := time.Second * time.Duration(n)
	fmt.Printf("server: call from client(%d) wait=%v \n", arg.Order, n)
	time.Sleep(wait)
	return nil
}

// -------------------------------------------------------------------
// test
// -------------------------------------------------------------------
func Test(t *testing.T) {
	const addr = "./test.sock"
	rand.Seed(time.Now().Unix())

	// 服务器部分
	//server, err := NewUnixServer(addr, nil)
	server, err := NewUnixServer(addr, jsonrpc.NewServerCodec)
	if err != nil {
		log.Fatal(err)
	}
	newRegister().RegisterTo(server, &S_Receiver{})
	go server.Run(context.Background())

	// 客户端部分
	//client, err := NewUnixClient(10, addr, nil)
	client, err := NewUnixClient(10, addr, jsonrpc.NewClientCodec)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("begin: len(idles)=%d\n", client.IdleCount())
	fun := func(wg *sync.WaitGroup, i int) {
		defer wg.Done()

		arg := S_Arg{Order: i, Message: fmt.Sprintf("msg: %d", i)}
		var ret S_Ret
		err := newRegister().CallWith(client, time.Second*5, &arg, &ret)
		if err != nil {
			fmt.Printf("clienti(%d): call fail, %v || len(idles)=%d\n", i, err, client.IdleCount())
		} else {
			fmt.Printf("client(%d): reply from server: %v || len(idles)=%d\n", i, ret, client.IdleCount())
		}
	}

	var wg sync.WaitGroup
	wg.Add(100)
	for i := range 100 {
		go fun(&wg, i)
		time.Sleep(time.Millisecond * 100)
	}

	wg.Wait()
	fmt.Printf("\nend: len(idles)=%d\n", client.IdleCount())

	// 退出
	client.Close()
	server.Close()
}
