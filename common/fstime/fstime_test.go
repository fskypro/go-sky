package fstime

import (
	"fmt"
	"testing"
	"time"

	"fsky.pro/fstest"
)

func TestFmt(t *testing.T) {
	fstest.PrintTestBegin("Fmt")
	defer fstest.PrintTestEnd()
}

func TestLastDayOfMon(t *testing.T) {
	fstest.PrintTestBegin("LastDayOfMon")
	defer fstest.PrintTestEnd()
	fmt.Println(LastDayOfMon(2022, time.February))
}

func TestDawn(t *testing.T) {
	fstest.PrintTestBegin("Dawn")
	defer fstest.PrintTestEnd()

	fmt.Println("Dawn(time.Now()): ", Dawn(time.Now()))
	fmt.Println(time.Now().Unix() - Dawn(time.Now()).Unix())

	fmt.Println("Dawn(time.Now().UTC()): ", Dawn(time.Now().UTC()))
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

func TestDaysBetween(t *testing.T) {
	t1 := time.Now()
	t2 := t1.AddDate(0, -1, 0)
	fmt.Println(DaysBetween(t1, t2))
}
