package fsnet

import "testing"
import "fmt"
import "fsky.pro/fstest"

func TestGetFreePort(t *testing.T) {
	fstest.PrintTestBegin("GetFreePort")

	port, err := GetFreePort()
	if err != nil {
		fmt.Print("Error: get free port fail:", err.Error())
	} else {
		fmt.Printf("get a free port: %d\n", port)
	}

	fstest.PrintTestEnd()
}

func TestInt64ToAddr(t *testing.T) {
	fstest.PrintTestBegin("Int64ToAddr/AddrToInt64")

	fmt.Printf(`call InetAtoN("192.168.1.2"): `)
	vaddr := InetAtoN("192.168.1.2")
	fmt.Println(vaddr)
	fmt.Printf("call InetNtoA(%d):  ", vaddr)
	fmt.Println(InetNtoA(vaddr))
	fmt.Println("")

	fmt.Printf(`call AddrToInt64("192.168.2.56", 0): `)
	vaddr = AddrToInt64("192.168.2.56", 0)
	fmt.Println(vaddr)
	fmt.Printf("call Int64ToAddr(%d):  ", vaddr)
	fmt.Println(Int64ToAddr(vaddr))
	fmt.Println("")

	fmt.Printf(`call AddrToInt64("192.168.2.57:12345", 0)`)
	vaddr = AddrToInt64("192.168.2.57:12345", 0)
	fmt.Println(vaddr)
	fmt.Printf("call Int64ToAddr(%d):  ", vaddr)
	fmt.Println(Int64ToAddr(vaddr))
	fmt.Println("")

	fmt.Printf(`call AddrToInt64("192.168.2.57", 12345):  `)
	vaddr = AddrToInt64("192.168.2.57", 12345)
	fmt.Println(vaddr)
	fmt.Printf("call Int64ToAddr(%d):  ", vaddr)
	fmt.Println(Int64ToAddr(vaddr))

	fstest.PrintTestEnd()
}
