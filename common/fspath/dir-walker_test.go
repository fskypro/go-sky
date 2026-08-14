package fspath

import (
	"fmt"
	"testing"

	"fsky.pro/fsky"
	"fsky.pro/fstest"
)

func TestDirWalker(t *testing.T) {
	fstest.PrintTestBegin("DirWalker")
	defer fstest.PrintTestEnd()

	walker := NewDirWalker("./")
	walker.Walk(func(fullPath string, info *S_FileInfo, err error) bool {
		fmt.Printf("path=%q, has-sub=%v\n", info.FullPath, fsky.Ret1[bool](info.HasChildren()))
		return true
	})
}
