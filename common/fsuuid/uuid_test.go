package fsuuid

import (
	"fmt"
	"sync"
	"testing"

	"fsky.pro/fstest"
)

func TestUUIDv1(t *testing.T) {
	fstest.PrintTestBegin("UUIDv1")
	defer fstest.PrintTestEnd()

	gen1 := func(wg *sync.WaitGroup) {
		uuid := NewV1()
		fmt.Println(uuid.String(), uuid.Version(), uuid.Variant())
		wg.Done()
	}
	gen4 := func(wg *sync.WaitGroup) {
		uuid := NewV4()
		fmt.Println(uuid.String(), uuid.Version(), uuid.Variant())
		wg.Done()
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(2)
		go gen1(&wg)
		go gen4(&wg)
	}
	wg.Wait()
}
