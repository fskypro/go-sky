package fsstrconv

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

func TestStrTo(t *testing.T) {
	fstest.PrintTestBegin("StrTo")
	defer fstest.PrintTestEnd()

	fmt.Println(StrTo[int]("32"))
	fmt.Println(StrTo[uint]("123"))
	fmt.Println(StrTo[bool]("ABC"))
	fmt.Println(StrTo[float32]("-122.3"))
	fmt.Println(StrTo[bool]("true"))
	fmt.Println(StrTo[float64]("88.888"))
}
