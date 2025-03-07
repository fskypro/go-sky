package fslua

import (
	"fmt"
	"log"
	"testing"

	"fsky.pro/fstest"
	lua "github.com/yuin/gopher-lua"
)

// 导出给 lua 的结构体父结构体
type Parent struct {
	Value int `lua:"value"`
}

// 只有公有方法才能导出给 lua 使用
func (self *Parent) Test() string {
	return "CCCC"
}

type Hand struct {
	name string `lua:"name"`
}

func (this *Hand) Func() (string, int) {
	return this.name, 50
}

// 传递给 lua 脚本的结构体
type S struct {
	*Parent
	name string `lua:"name"`
	Age  int
	Hand *Hand `lua:"hand"`
}

func (this *S) Func() (string, int) {
	return this.name, this.Age
}

func (this *S) SetAge(age int) {
	this.Age = age
}

func TestOO(t *testing.T) {
	fstest.PrintTestBegin("oo")
	defer fstest.PrintTestEnd()

	L := NewState()
	defer L.Close()

	file := "./scripts/oo.lua"
	if err := L.DoFile(file); err != nil {
		log.Fatalf("load lua file %q fail, %v", file, err)
	}

	// go 对象
	s := &S{
		Parent: &Parent{500},
		name:   "fanky",
		Age:    40,
		Hand:   &Hand{"right hand"},
	}
	// 封装 go 对象
	ud, err := WrapUserData(L, s)
	if err != nil {
		log.Fatalf("wrap s fail, %v", err)
	}

	// 调用 lua 脚本中的对象函数，并且将 go 对象作为参数传递给 lua 函数
	// 获取 lua 脚本中的全局对象
	obj1 := L.GetGlobal("obj1")
	if obj1 == lua.LNil {
		log.Fatal("obj1 not found in Lua script")
	}

	str, err := L.CallMethod1(obj1, "printInfo1", ud)
	if err != nil {
		log.Fatal("call printInfo fail,", err)
	}
	fmt.Println("return 1:", str)

	obj2 := L.GetGlobal("obj2")
	if obj2 == lua.LNil {
		log.Fatal("obj2 not found in Lua script")
	}

	str, err = L.CallMethod1(obj2, "printInfo2", ud)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("return 2:", str)

	fmt.Println(s.Age)
}
