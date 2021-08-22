package fstime

import (
	"fmt"
	"testing"
	"time"

	"fsky.pro/fstest"
)

func TestDawn(t *testing.T) {
	fstest.PrintTestBegin("Dawn")
	fmt.Println("Dawn(time.Now()): ", Dawn(time.Now()))

	fmt.Println("Dawn(time.Now().UTC()): ", Dawn(time.Now().UTC()))
	fstest.PrintTestEnd()
}

func TestWeekStart(t *testing.T) {
	fstest.PrintTestBegin("WeekStart")
	defer fstest.PrintTestEnd()

	fmt.Println("WeekStart(time.Now()) of WeekStartMonday: ", WeekStart(time.Now()))

	CWeekStart = WeekStartSunday
	fmt.Println("WeekStart(time.Now()) of WeekStartSunday: ", WeekStart(time.Now()))

	fmt.Println("WeekStart(time.UTC()) of WeekStartSunday: ", WeekStart(time.Now().UTC()))

	d, h, m, s := Seconds2DaysTime(3600*24 + 100)
	fmt.Printf("%d天%d小时%d分钟%d秒\n", d, h, m, s)
}
