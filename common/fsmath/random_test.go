package fsmath

import "fmt"
import "testing"
import "fsky.pro/fstest"

func TestRandSecond(t *testing.T) {
	fstest.PrintTestBegin("RandSecond")

	for i := 0; i < 10; i += 1 {
		fmt.Println(RandSecond(50, 101))
	}

	fstest.PrintTestEnd()
}

func TestRandMillisecond(t *testing.T) {
	fstest.PrintTestBegin("RandMillisecond")

	for i := 0; i < 10; i += 1 {
		fmt.Println(RandMillisecond(50, 101))
	}

	fstest.PrintTestEnd()
}
