package fstype

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

func TestStrToNumber(t *testing.T) {
	fstest.PrintTestBegin("StrToNumber")
	defer fstest.PrintTestEnd()

	fmt.Println(StrToNumber[int]("100"))
}
