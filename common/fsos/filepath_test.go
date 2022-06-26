package fsos

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

func TestIsPathExists(t *testing.T) {
	fstest.PrintTestBegin("IsPathExists")

	fmt.Println(IsPathExists("/data/webroot"))
	fmt.Println(IsPathExists("xxx"))

	fstest.PrintTestEnd()
}

func TestIsDirExists(t *testing.T) {
	fstest.PrintTestBegin("IsDirExists")

	fmt.Println(IsDirExists("/data/webroot"))
	fmt.Println(IsDirExists("~/.bashrc"))

	fstest.PrintTestEnd()
}

func TestIsFileExists(t *testing.T) {
	fstest.PrintTestBegin("IsFileExists")

	fmt.Println(IsFileExists("/data/webroot"))
	fmt.Println(IsFileExists("/root/.bashrc"))

	fstest.PrintTestEnd()
}

func TestGetFullPathToBin(t *testing.T) {
	fstest.PrintTestBegin("GetFullPathToBin")
	fmt.Println(GetFullPathToBin("/abc/def/ghijk.bin"))
	fmt.Println(GetFullPathToBin("./abc/def/ghijk.bin"))
	fstest.PrintTestEnd()
}

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
