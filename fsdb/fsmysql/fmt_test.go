package fsmysql

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

func TestFmtSelectPrepare(t *testing.T) {
	fstest.PrintTestBegin("FmtSelectPrepare")

	var aaa string
	var bbb int
	var ccc float32

	cols := map[string]interface{}{
		"aaa": &aaa,
		"bbb": &bbb,
		"ccc": &ccc,
	}

	sqltx, _ := FmtSelectPrepare(cols, "`table`", "")
	fmt.Println(sqltx)

	fstest.PrintTestEnd()
}

func TestFmtInsertPrepare(t *testing.T) {
	fstest.PrintTestBegin("FmtInsertPrepare")

	cols := map[string]interface{}{
		"aaa": "aaa",
		"bbb": []byte("bbb"),
		"ccc": 1000,
	}

	sqltx, _ := FmtInsertPrepare("table", cols)
	fmt.Println(sqltx)

	fstest.PrintTestEnd()
}

func TestFmtUpdatePrepare(t *testing.T) {
	fstest.PrintTestBegin("FmtUpdatePrepare")

	cols := map[string]interface{}{
		"aaa": "aaa",
		"bbb": []byte("bbb"),
		"ccc": 1000,
		"ddd": Unquote{"`ddd`+300"},
	}

	sqltx, _ := FmtUpdatePrepare("table", cols, "where `aaa`='aaa'")
	fmt.Println(sqltx)

	fstest.PrintTestEnd()
}
