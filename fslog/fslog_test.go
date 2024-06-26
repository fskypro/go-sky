package fslog

import (
	"testing"

	"fsky.pro/fstest"
)

func TestDebug(t *testing.T) {
	fstest.PrintTestBegin("Debug")

	// fslog
	Debug("test debug 111111!")
	Debugf("test %s 222222!", "debug")

	Shield("debug")
	Debug("test debug 333333!")

	Unshield("Debug")
	Debug("test debug 444444!")

	Unshield("info", "error")
	Debug("test debug 555555!")

	Unshield("info", "error", "debug", "hack", "trace")
	Debug("test debug 666666!")

	// fslogger
	fl := NewDayfileLogger("./logs", "test")
	SetLogger(fl)
	Debug("debug", "aaaaaaa")
	fstest.PrintTestEnd()
}

func TestInfo(t *testing.T) {
	fstest.PrintTestBegin("Info")

	// fslog
	Info("test info 111111!")
	Infof("test %s 222222!", "info")

	Shield("info")
	Info("test info 333333!")

	Unshield("info")
	Info("test info 444444!")

	// fslogger
	fl := NewDayfileLogger("./logs", "test")
	SetLogger(fl)
	Info("info", "aaaaaaa")
	fl.SetNewLogCmd("./linklog.sh", "arg")
	fstest.PrintTestEnd()
}

func TestError(t *testing.T) {
	fstest.PrintTestBegin("Error")

	// fslog
	Error("test error 111111!")
	Errorf("test %s 222222!", "error")

	Shield("error")
	Error("test error 333333!")

	Shield("error")
	Error("test error 444444!")

	// fslogger
	fl := NewDayfileLogger("./logs", "test")
	SetLogger(fl)
	Error("error", "aaaaaaa")

	fstest.PrintTestEnd()
}

func TestHack(t *testing.T) {
	fstest.PrintTestBegin("Hack")

	// fslog
	Hack("test hack 111111!")
	Hackf("test %s 222222!", "hack")

	Shield("hack")
	Hack("test hack 333333!")

	Unshield("hack")
	Hack("test hack 444444!")

	// fslogger
	fl := NewDayfileLogger("./logs", "test")
	SetLogger(fl)
	Hack("hack", "aaaaaaa")

	fstest.PrintTestEnd()
}

func TestIlleg(t *testing.T) {
	fstest.PrintTestBegin("Illeg")

	// fslog
	Illeg("test hack 111111!")
	Illegf("test %s 222222!", "illeg")

	Shield("hack")
	Illeg("test hack 333333!")

	Unshield("hack")
	Illeg("test hack 444444!")

	// fslogger
	fl := NewDayfileLogger("./logs", "test")
	SetLogger(fl)
	Illeg("hack", "aaaaaaa")

	fstest.PrintTestEnd()
}

func TestTrace(t *testing.T) {
	fstest.PrintTestBegin("Trace")

	Trace("test trace 111111!")
	Tracef("test %s 222222!", "trace")

	Shield("trace")
	Trace("test trace 333333!")

	Unshield("trace")
	Trace("test trace 444444!")

	// fslogger
	fl := NewDayfileLogger("./logs", "test")
	// 替换全局 fslogger
	SetLogger(fl)
	Trace("trace", "aaaaaaa") // 直接用 fslog.Trace 即可

	fstest.PrintTestEnd()
}

func TestFatal(t *testing.T) {
	Fatal("xxxxxxx")
}
