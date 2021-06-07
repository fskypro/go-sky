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
		"bbb": 10.1,
		"ccc": 1000,
	}

	sqltx, values := FmtInsertPrepare("table", cols)
	fmt.Println(sqltx)
	for i, v := range values {
		fmt.Printf("values[%d]=%v\n", i, v)
	}

	fstest.PrintTestEnd()
}

func TestFmtInsertsPrepare(t *testing.T) {
	fstest.PrintTestBegin("FmtInsertsPrepare")

	vs1 := map[string]interface{}{
		"aaa": "aaa",
		"bbb": 1.1,
		"ccc": 1000,
	}
	vs2 := map[string]interface{}{
		//"aaa": "xxxx",
		"bbb": 2.1,
		"ccc": 2000,
	}

	sqltx, values := FmtInsertPrepare("table", vs1, vs2)
	fmt.Println(sqltx)
	for i, v := range values {
		fmt.Printf("values[%d]=%v\n", i, v)
	}

	fstest.PrintTestEnd()
}

func TestFmtInsertIgnorePrepare(t *testing.T) {
	fstest.PrintTestBegin("FmtInsertPrepare")

	cols := map[string]interface{}{
		"aaa": "aaa",
		"bbb": []byte("bbb"),
		"ccc": 1000,
	}

	sqltx, _ := FmtInsertIgnorePrepare("table", cols)
	fmt.Println(sqltx)

	fstest.PrintTestEnd()
}

func TestFmtInsertUpdatePrepare(t *testing.T) {
	fstest.PrintTestBegin("FmtInsertPrepare")

	inserts := map[string]interface{}{
		"aaa": "aaa",
		"bbb": []byte("bbb"),
		"ccc": 1000,
	}

	updates := map[string]interface{}{
		"aaa": "bbb",
		"bbb": []byte("ccc"),
		"ccc": 2000,
	}

	sqltx, _ := FmtInsertUpdatePrepare("table", inserts, updates)
	fmt.Println(sqltx)

	fstest.PrintTestEnd()
}

func TestFmtUpdatePrepare(t *testing.T) {
	fstest.PrintTestBegin("FmtUpdatePrepare")

	cols := map[string]interface{}{
		"aaa": "aaa",
		"bbb": []byte("bbb"),
		"ccc": 1000,
		"ddd": "`ddd`+300",
	}

	sqltx, values := FmtUpdatePrepare("table", cols, "`aaa`=?", "abc")
	fmt.Println(values)
	fmt.Println(sqltx)

	fstest.PrintTestEnd()
}
