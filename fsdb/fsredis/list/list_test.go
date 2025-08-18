package list

import (
	"fmt"
	"log"
	"testing"

	"fsky.pro/fstest"
	"github.com/redis/go-redis/v9"
)

type S_Item struct {
	StrValue string
	IntValue int
}

func TestList(t *testing.T) {
	fstest.PrintTestBegin("AnyKeySet")
	defer fstest.PrintTestEnd()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "@yxht123456YXHT",
		DB:       0,
	})
	a1 := &S_Item{"xxxx", 100}
	a2 := &S_Item{"yyyy", 200}
	a3 := &S_Item{"zzzz", 300}
	a4 := &S_Item{"wwww", 400}

	list := NewList[*S_Item](rdb, "test-list", nil)
	list.Clear()
	item, err := list.LPop()
	if err != nil {
		log.Fatal(1111, err)
	}
	fmt.Println("aaaaa", item)

	n1, err := list.RPush(a1)
	if err != nil {
		log.Fatal(22222, err)
	}
	n2, err := list.RPush(a2, a3)
	if err != nil {
		log.Fatal(3333, err)
	}
	fmt.Println("bbbbb:", n1, n2)

	// ------------------------------
	err = list.Fetch(func(item S_RetValue[*S_Item]) bool {
		fmt.Printf("ccccc:%#v\n", item.Value)
		return true
	})
	if err != nil {
		log.Fatal(44444, err)
	}

	// ------------------------------
	items, err := list.GetRange(0, 1)
	if err != nil {
		log.Fatal("5555", err)
	}
	for _, item := range items {
		fmt.Printf("dddddd: %#v\n", item.Value)
	}

	// ------------------------------
	n, err := list.InsertBefore(a3, a4)
	if err != nil {
		log.Fatal(55555, err)
	}
	fmt.Println("ddddd", n)

	items, err = list.GetRange(0, -1)
	if err != nil {
		log.Fatal("66666", err)
	}
	for _, item := range items {
		fmt.Printf("ffff: %#v\n", item.Value)
	}

	// ------------------------------
	item, err = list.RPop()
	if err != nil {
		log.Fatal(77777, err)
	}
	fmt.Println("gggggg", item.Exists, item.Value)

	// ------------------------------
	item, err = list.Get(-1)
	if err != nil {
		log.Fatal(88888, err)
	}
	fmt.Println("hhhhh", item.Exists, item.Value)

	// ------------------------------
	ok, err := list.RemoveValue(a3)
	if err != nil {
		log.Fatal(999999)
	}
	fmt.Println("iiiiiii", ok)

	// ------------------------------
	ok, err = list.Remove(1)
	if err != nil {
		log.Fatal("AAAAAA", err)
	}
	fmt.Println("jjjjj", ok)

	// ------------------------------
	err = list.Fetch(func(item S_RetValue[*S_Item]) bool {
		fmt.Println("kkkkkkk", item.Value)
		return true
	})
	if err != nil {
		log.Fatal("BBBBB", err)
	}
}
