package fsio

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

func TestCopyFile(t *testing.T) {
	fstest.PrintTestBegin("CopyFile")
	src, dst := "./util.go", "util2.txt"
	err := CopyFile(src, dst)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("copyfile success: %q -> %q\n", src, dst)
	}
	fstest.PrintTestEnd()
}
