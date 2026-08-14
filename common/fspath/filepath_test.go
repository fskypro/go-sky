package fspath

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	"fsky.pro/fstest"
)

func TestFileNoExt(t *testing.T) {
	fstest.PrintTestBegin("FileNoExt")
	defer fstest.PrintTestEnd()
	fmt.Println(FileNoExt("aaa/bbb/ccc.txt"))
	fmt.Println(FileNoExt("aaa/bbb/ccc"))
}

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

func TestWorkDir(t *testing.T) {
	fstest.PrintTestBegin("WorkDir")
	defer fstest.PrintTestEnd()
	WorkDir("."+string(os.PathSeparator), func(file string, info fs.FileInfo) bool {
		fmt.Println(file, info)
		return true
	})
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
