package fsreflect

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

type item struct {
	count int
}

type mlist struct {
	list [][]*item
}

type user struct {
	Map *mlist
}

func TestDeepFieldValue(t *testing.T) {
	fstest.PrintTestBegin("Get/SetDeepFieldValue")
	defer fstest.PrintTestEnd()

	u := new(user)
	u.Map = &mlist{
		list: [][]*item{[]*item{new(item)}},
	}

	v, er := GetDeepFieldValue(u, "Map.list[0][0].count")
	if er != nil {
		fmt.Println("deep get struct field value fail:", er.Error())
	} else {
		fmt.Println("user.Map.list[0][0].count =", v)
	}

	if er = SetDeepFieldValue(u, "Map.list[0][0].count", 100); er != nil {
		fmt.Println("deep set struct field value fail:", er.Error())
	} else {
		fmt.Println("deep set struct field value success, new value =", u.Map.list[0][0].count)
	}
}

func TestCopyObject(t *testing.T) {
	fstest.PrintTestBegin("CopyStructObject")
	defer fstest.PrintTestEnd()

	type Object struct {
		xxx int
		yyy string
		zzz float32
	}

	src := &Object{10, "abcd", 23.3}
	dst, err := CopyStructObject(src)
	if err != nil {
		fmt.Println("copy object fail:", err.Error())
		return
	}
	fmt.Printf("copy object success, dst = %v\n", dst)
}
