package fstime

import (
	"fmt"
	"testing"
	"time"

	"fsky.pro/fstest"
)

func TestTimer(t *testing.T) {
	fstest.PrintTestBegin("Timer")
	defer fstest.PrintTestEnd()

	timer := NewTimer(time.Second * 10)
	go func() {
		a := <-timer.C
		fmt.Println(11111, a)
	}()

	time.Sleep(5 * time.Second)
	timer.Cancel()

	time.Sleep(11 * time.Second)
}
