package fslua

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"fsky.pro/fsstr/fsfmt"
	"fsky.pro/fstest"
	glua "github.com/yuin/gopher-lua"
	lua "github.com/yuin/gopher-lua"
)

// -------------------------------------------------------------------
// 自定义解释函数
// -------------------------------------------------------------------
type I_MyTime interface {
	FmtString() string
}

type MyTime time.Time

func (this *MyTime) FmtString() string {
	return (time.Time)(*this).Format(time.DateTime)
}

func (this *MyTime) UnmarshalLua(lv glua.LValue) error {
	t, err := time.ParseInLocation(time.DateTime, lv.String(), time.Local)
	if err != nil {
		return fmt.Errorf("can't parse time string %q to format %q", lv.String(), time.DateTime)
	}
	*this = MyTime(t)
	return nil
}

// -------------------------------------------------------------------
// 基类结构
// -------------------------------------------------------------------
type s_Base struct {
	BaseName string `lua:"BaseName"`
	Array    []int
}

type s_Base2 struct {
	Base2Name string `lua:"Base2Name"`
	Slice     []any
}

// -------------------------------------------------------------------
// 嵌套结构
// -------------------------------------------------------------------
type S_NestStruct struct {
	A string `lua:"a"`
	B int    `lua:"b"`
}

type s_SubConfig struct {
	Name      string `lua:"name"`
	IntValue  int8   `lua:"intValue"`
	boolValue bool

	NestStructs *map[string]*S_NestStruct
}

// -------------------------------------------------------------------
// 主表
// -------------------------------------------------------------------
type s_Config struct {
	*s_Base
	s_Base2
	Unexposed string `lua:"-"`

	Name       string
	IntValue   int8
	FloatValue float32
	BoolValue  bool
	Sub        s_SubConfig `lua:"subConfig"`

	inner *struct {
		Map    map[int]string `lua:"map"`
		MyTime *MyTime        `lua:"myTime" fsfmt:"str"`
		//IMyTime I_MyTime       `lua:"myTime" fsfmt:"str"`
	} `lua:"inner"`

	NestSice []map[string]int `lua:"nestSlice"`
	NestAnys map[string]any
}

// -------------------------------------------------------------------
// 创建 lua 虚拟机
// -------------------------------------------------------------------
func newLuaState() *glua.LState {
	L := glua.NewState()
	err := L.DoFile("./scripts/table.lua")
	if err != nil {
		log.Fatalf("load lua fail, %v", err)
	}
	return L
}

// -------------------------------------------------------------------
// 测试函数
// -------------------------------------------------------------------
// lua table 解释为 go 对象
func TestUnmarshalTable(t *testing.T) {
	fstest.PrintTestBegin("unmarshal lua table")
	defer fstest.PrintTestEnd()
	L := newLuaState()
	defer L.Close()

	g := L.GetGlobal("config")
	tb := g.(*glua.LTable)

	// 解释为结构体
	conf := new(s_Config)
	err := UnmarshalTable(tb, conf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fsfmt.SprintStruct(conf, nil))
	//fmt.Println(conf.inner.IMyTime)
	fmt.Println("-------------------")

	// 解释为笼统的 Map
	mapConf := map[string]any{}
	err = UnmarshalTable(tb, &mapConf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fsfmt.SprintStruct(mapConf, nil))

	fmt.Println("-------------------")
	jdata, err := json.MarshalIndent(mapConf, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jdata))
}

type I interface {
	Test()
}

type A struct {
	Str  string
	List []int
	Map  *map[int]string
}

func (*A) Test() {}

type B struct {
	A *A
}

// go 对象解释为 lua table
func TestMarshalTable(t *testing.T) {
	fstest.PrintTestBegin("marshal lua table")
	defer fstest.PrintTestEnd()

	L := newLuaState()
	defer L.Close()

	b := &B{
		A: &A{
			Str:  "xxx",
			List: []int{100, 200, 300},
			Map:  &map[int]string{1: "xxx", 2: "yy"},
		},
	}
	tb, err := MarshalTable(b)
	if err != nil {
		log.Fatal(err)
	}

	fnPrintTable := L.GetGlobal("printTable")
	if fnPrintTable.Type() == glua.LTFunction {
		err := L.CallByParam(lua.P{
			Fn:      fnPrintTable,
			NRet:    1,
			Protect: true,
		}, tb)
		if err != nil {
			log.Fatalf("call lua function fail, %v\n", err)
		}
	}

	fmt.Println("-------------------")
	bb := new(B)
	if err = UnmarshalTable(tb, bb); err != nil {
		log.Fatal(err)
	}
	fmt.Println(fsfmt.SprintStruct(bb, nil))
}
