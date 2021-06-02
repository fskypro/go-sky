package fsbytes

import "fmt"
import "testing"
import "fsky.pro/fstest"

func TestIndexN(t *testing.T) {
	fstest.PrintTestBegin("IndexN")
	fmt.Println(`IndexN("123", 0, "123"): `, IndexN([]byte("123"), 0, []byte("123")))
	fmt.Println(`IndexN("123", 1, "123"): `, IndexN([]byte("123"), 1, []byte("123")))
	fmt.Println(`IndexN("13", 0, "123"): `, IndexN([]byte("13"), 0, []byte("123")))
	fmt.Println(`IndexN("124123", 1, "123"): `, IndexN([]byte("124123"), 1, []byte("123")))
	fstest.PrintTestEnd()
}
