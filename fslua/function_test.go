package fslua

import (
	"fmt"
	"log"
	"testing"

	"fsky.pro/fslog"
	"fsky.pro/fsstr/fsfmt"
	"fsky.pro/fstest"
	lua "github.com/yuin/gopher-lua"
)

type s_Ret struct {
	StrValue   string  `lua:"strValue"`
	FloatValue float32 `lua:"floatValue"`
}

func TestFunction(t *testing.T) {
	fstest.PrintTestBegin("test function")
	defer fstest.PrintTestEnd()

	L := lua.NewState()
	err := L.DoFile("./scripts/function.lua")
	if err != nil {
		log.Fatalf("load lua fail, %v", err)
	}
	defer L.Close()

	SetGlobalFunc(L, "fs_debugf", fslog.Debugf)

	a, b, r, err := CallGlobalFunc3[int, string, *s_Ret](L, "func", "lua")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("returns: (%d, %q)\n%s\n", a, b, fsfmt.SprintStruct(r, nil))
}
