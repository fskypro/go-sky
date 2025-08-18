package hset

import (
	"fmt"
	"log"
	"testing"
	"time"

	"fsky.pro/fsredis/hset"
	"fsky.pro/fstest"

	"github.com/redis/go-redis/v9"
)

type S_AnyValue struct {
	Int int
	Str string
	Map map[string]int
}

type S_AnyKey struct {
	text string
}

func (self S_AnyKey) KeyToString() string {
	return self.text
}

func (this S_AnyKey) StringToKey(str string) (S_AnyKey, error) {
	return S_AnyKey{text: str}, nil
}

func TestAnyKeySet(t *testing.T) {
	fstest.PrintTestBegin("AnyKeySet")
	defer fstest.PrintTestEnd()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "@yxht123456YXHT",
		DB:       0,
	})

	m := map[S_AnyKey]*S_AnyValue{
		S_AnyKey{"a1"}: &S_AnyValue{Int: 1, Str: "hello", Map: map[string]int{"x": 1}},
		S_AnyKey{"a2"}: &S_AnyValue{Int: 2, Str: "world", Map: map[string]int{"y": 2}},
	}

	// 把 m 的每个元素作为一个 field 存入 redis hash
	hashKey := "str_key_map"
	s := hset.NewAnyKeySet[S_AnyKey, *S_AnyValue](rdb, hashKey, nil)
	err := s.Update(m)
	if err != nil {
		log.Fatal("aaaaa", err)
	}

	time.Sleep(time.Second * 1)

	fmt.Println(s.Get(S_AnyKey{"xx"}))

	err = s.Fetch("", func(item hset.S_IterItem[S_AnyKey, *S_AnyValue]) bool {
		fmt.Println(item.Key, ":", item.Value)
		return true
	})
	if err != nil {
		log.Fatal("error:", err)
	}
}
