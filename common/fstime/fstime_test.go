package fstime

import "fmt"
import "time"
import "testing"
import "fsky.pro/fstest"

func TestDawn(t *testing.T) {
	fstest.PrintTestBegin("Dawn")
	fmt.Println("Dawn(time.Now()): ", Dawn(time.Now()))

	fmt.Println("Dawn(time.Now().UTC()): ", Dawn(time.Now().UTC()))
	fstest.PrintTestEnd()
}

func TestWeekStart(t *testing.T) {
	fstest.PrintTestBegin("WeekStart")
	fmt.Println("WeekStart(time.Now()) of WeekStartMonday: ", WeekStart(time.Now()))

	CWeekStart = WeekStartSunday
	fmt.Println("WeekStart(time.Now()) of WeekStartSunday: ", WeekStart(time.Now()))

	fmt.Println("WeekStart(time.UTC()) of WeekStartSunday: ", WeekStart(time.Now().UTC()))
	fstest.PrintTestEnd()
}
