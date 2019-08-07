package fsutil

import "testing"
import "fmt"
import "fsky.pro/fstest"

func TestUnzip(t *testing.T) {
	fstest.PrintTestBegin("Unzip/Zip")

	var s = "世界你好"
	data := []byte(s)
	fmt.Println("unzip data: ", data)

	zdata, err := Zip(data)
	if err != nil {
		fmt.Println("zip error: ", err.Error())
		return
	}
	fmt.Println("zip data: ", zdata)

	data, err = Unzip(zdata)
	if err != nil {
		fmt.Println("unzip error: ", err.Error())
		return
	}
	fmt.Println("unzip data: ", data)

	fstest.PrintTestEnd()
}
