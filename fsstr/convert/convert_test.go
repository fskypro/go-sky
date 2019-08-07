package convert

import "testing"
import "fmt"
import "fsky.pro/fstest"

func TestBytes2String(t *testing.T) {
	fstest.PrintTestBegin("Bytes2String")

	var bstr = []byte("112233")
	var str = Bytes2String(bstr)
	fmt.Println(str, len(str))
	bstr[1] = 'e'
	fmt.Println(str)
	fstest.PrintTestEnd()
}

func TestString2Bytes(t *testing.T) {
	fstest.PrintTestBegin("String2Bytes")

	str := "xxxxxyyyyy"
	var bstr = String2Bytes(str)
	fmt.Println(bstr)

	fstest.PrintTestEnd()
}
