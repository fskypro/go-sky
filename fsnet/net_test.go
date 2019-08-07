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
