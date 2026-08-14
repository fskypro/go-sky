package hset

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"fsky.pro/fsredis/hset"
	"fsky.pro/fstest"

	"github.com/redis/go-redis/v9"
)

type S_StrValue struct {
	Int int
	Str string
	Map map[string]int
}

type T_StrKey string

func TestStrKeySet(t *testing.T) {
	fstest.PrintTestBegin("StrKeySet")
	defer fstest.PrintTestEnd()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "@yxht123456YXHT",
		DB:       0,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	st := rdb.Ping(ctx)
	if st.Err() != nil {
		log.Fatalf("ping redis db fail, %v", st.Err())
	}

	m := map[T_StrKey]*S_StrValue{
		T_StrKey("a1"): &S_StrValue{Int: 1, Str: "hello", Map: map[string]int{"x": 1}},
		T_StrKey("a2"): &S_StrValue{Int: 2, Str: "world", Map: map[string]int{"y": 2}},
	}

	// 把 m 的每个元素作为一个 field 存入 redis hash
	hashKey := "str_key_map"
	s := hset.NewStrKeySet[T_StrKey, *S_StrValue](rdb, hashKey, nil)
	err := s.Update(m)
	if err != nil {
		log.Fatal("aaaaa", err)
	}

	time.Sleep(time.Second * 1)

	fmt.Println(s.Get("xx"))

	err = s.Fetch("", func(item hset.S_IterItem[T_StrKey, *S_StrValue]) bool {
		fmt.Println(item.Key, ":", item.Value)
		return true
	})
	if err != nil {
		log.Fatal("error:", err)
	}
}

func TestExpire(t *testing.T) {
	fstest.PrintTestBegin("Expire")
	defer fstest.PrintTestEnd()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "@yxht123456YXHT",
		DB:       0,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	st := rdb.Ping(ctx)
	if st.Err() != nil {
		log.Fatalf("ping redis db fail, %v", st.Err())
	}

	hashKey := "test_expire"
	s := NewStrKeySet[string, int](rdb, hashKey, nil)
	err := s.Set("xxx", 100)
	if err != nil {
		log.Fatal("aaaaa", err)
	}

	ret, err := s.Get("xxx")
	if err != nil {
		log.Fatalf("get value of key=%v fail, %v", "xxx", err)
	}
	if ret.Exists {
		fmt.Printf("value of key=%v is %v\n", "xxx", ret.Value)
	}

	s.client.HExpire(context.Background(), hashKey, time.Second*2, "xxx")

	time.Sleep(time.Second * 3)

	ret, err = s.Get("xxx")
	if err != nil {
		log.Fatalf("get value of key=%v fail, %v", "xxx", err)
	}
	if !ret.Exists {
		fmt.Printf("value of key=%v has not exists\n", "xxx")
	}
}
