package fsstrconv

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

func TestStr2TypeOf(t *testing.T) {
	fstest.PrintTestBegin("Str2TypeOf")
	defer fstest.PrintTestEnd()

	fmt.Println(Str2TypeOf("", 32))
	fmt.Println(Str2TypeOf("123", 32))
	fmt.Println(Str2TypeOf("ABC", []byte{}))
	fmt.Println(Str2TypeOf("sf", nil))
	fmt.Println(Str2TypeOf("-122.3", 32))
	fmt.Println(Str2TypeOf("true", false))
	fmt.Println(Str2TypeOf("88.888", 0.1))
	var a *int = nil
	fmt.Println(Str2TypeOf("23", a))
}
