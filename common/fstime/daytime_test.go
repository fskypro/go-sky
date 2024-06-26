package fstime

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

func TestDayTime(t *testing.T) {
	fstest.PrintTestBegin("DayTime")
	defer fstest.PrintTestEnd()

	dt := NewDayTime(0, 0, 59)
	fmt.Println(dt.Add(-1, 0, 0))

	dt = NewDayTime(24,0,0)
	fmt.Println(dt)

	fmt.Println(ParseDayTime(""))
}
