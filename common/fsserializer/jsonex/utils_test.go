package jsonex

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"fsky.pro/fstest"
)

type Obj struct {
	StrValue string `json:"strValue"`
	IntValue int    `json:"intValue"`
}

func TestObjectJsonMap(t *testing.T) {
	fstest.PrintTestBegin("ObjectJsonMap")
	defer fstest.PrintTestEnd()
	obj := &Obj{"xxx", 100}
	m, err := ObjectJsonMap(obj)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", m)
}

// -------------------------------------------------------------------
// TestMarshalJson
// -------------------------------------------------------------------
type A struct {
	T time.Time `json:"t" out:"tt"`
}

type B struct {
	*A
	Int int    `json:"int"`
	Str string `out:"sss"`
}

type C struct {
	*A
	Int int    `json:"int"`
	Str string `out:"sss"`
}

func (self C) MarshalJSON() ([]byte, error) {
	return MarshalJsonWithTag(&self, "out")
}

func TestMarshalJsonWithTag(t *testing.T) {
	fstest.PrintTestBegin("MarshalJSON")
	defer fstest.PrintTestEnd()
	// 直接调用 MarshalJsonWithTag
	b := &B{
		A: &A{
			T: time.Now(),
		},
		Int: 100,
		Str: "str in B",
	}
	data, err := MarshalJsonWithTag(b, "out")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))

	// 通过 MarshalJSON 函数转调用 MarshalJsonWithTag
	c := &C{
		A: &A{
			T: time.Now(),
		},
		Int: 100,
		Str: "str in B",
	}
	data, err = json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))

}
