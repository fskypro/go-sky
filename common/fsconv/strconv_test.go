package fsconv

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

type Str string
type Int8 int8
type Float32 float32

func TestStrTo(t *testing.T) {
	fstest.PrintTestBegin("StrTo")
	defer fstest.PrintTestEnd()

	fmt.Println(StrTo[int]("32x"))
	fmt.Println(StrTo[uint]("123"))
	fmt.Println(StrTo[bool]("ABC"))
	fmt.Println(StrTo[float32]("-122.3"))
	fmt.Println(StrTo[bool]("true"))
	fmt.Println(StrTo[float64]("88.888"))

	fmt.Println(StrTo[Str]("xxxx"))
	fmt.Println(StrTo[Int8]("3200"))
	fmt.Println(StrTo[Float32]("500.6"))
}
