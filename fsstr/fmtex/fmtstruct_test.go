package fmtex

import (
	"fmt"
	"testing"
	"unsafe"

	"fsky.pro/fstest"
)

type s_NetInfo struct {
	Host string
	Port uint16
}

type S_DBInfo struct {
	Host string
	Port uint16
}

type AA struct {
	v int

	netInfo s_NetInfo
	dbInfo  *S_DBInfo
}

type Config struct {
	AA

	_bool      bool
	_pbool     *bool
	_int       int
	_pint      *int
	_int16     int16
	_pint16    *int16
	_int64     int64
	_pint64    *int64
	_uint      uint
	_puint     *uint
	_uint32    uint32
	_puint32   *uint32
	_uintptr   uintptr
	_float32   float32
	_pfloat64  *float64
	_complex64 complex64
	_pcomp128  *complex128
	_string    string
	_pstring   *string
	_up        unsafe.Pointer

	_array   [2]int
	_parray  *[2]string
	_dbarray [1]S_DBInfo

	_slice    []int
	_pslice   *[]string
	_dbslice  []*S_DBInfo
	_pdbslice *[]*S_DBInfo

	_map    map[int]string
	_pmap   *map[int]string
	_dbmap  map[int]S_DBInfo
	_pdbmap *map[int]*S_DBInfo

	_struct    s_NetInfo
	_pstruct   *s_NetInfo
	_nilstruct *s_NetInfo

	_chan  chan int
	_pchan *chan string

	_func      func(int, string) error
	_interface interface{}

	_nest *AA
}

func TestSprintStruct(t *testing.T) {
	fstest.PrintTestBegin("SprintStruct")
	_pbool := true
	_pint := -200
	var _pint16 int16 = -400
	var _pint64 int64 = -600
	var _puint uint = 200
	var _puint32 uint32 = 400
	_pfloat64 := 2.5
	_pcomplex128 := new(complex128)
	_pstring := "yyyy"

	_array := [2]int{1, 2}
	_parray := [2]string{"123", "456"}
	_dbarray := [1]S_DBInfo{S_DBInfo{"host", 80}}

	_pslice := []string{"aaaa", "bbbb", "cccc"}
	_dbslice := []*S_DBInfo{&S_DBInfo{}, &S_DBInfo{"host", 8080}}

	_map := map[int]string{1: "xxx", 2: "yyy"}
	_dbmap := map[int]S_DBInfo{1: S_DBInfo{"host", 1000}}
	_pdbmap := map[int]*S_DBInfo{2: &S_DBInfo{"host1", 2000}, 3: &S_DBInfo{"host2", 3000}}

	_struct := s_NetInfo{"net host", 123}
	_pstruct := &s_NetInfo{"net host", 456}

	_pchan := make(chan string)

	_func := func(int, string) error { return nil }

	cfg := &Config{
		AA: AA{123456, s_NetInfo{}, &S_DBInfo{"host", 123}},

		_bool:      false,
		_pbool:     &_pbool,
		_int:       -100,
		_pint:      &_pint,
		_int16:     -300,
		_pint16:    &_pint16,
		_int64:     -500,
		_pint64:    &_pint64,
		_uint:      100,
		_puint:     &_puint,
		_uint32:    00,
		_puint32:   &_puint32,
		_float32:   -1.5,
		_pfloat64:  &_pfloat64,
		_complex64: *new(complex64),
		_pcomp128:  _pcomplex128,
		_string:    "xxxxx",
		_pstring:   &_pstring,

		_array:   _array,
		_parray:  &_parray,
		_dbarray: _dbarray,

		_slice:    nil,
		_pslice:   &_pslice,
		_dbslice:  _dbslice,
		_pdbslice: &_dbslice,

		_map:    _map,
		_pmap:   &_map,
		_dbmap:  _dbmap,
		_pdbmap: &_pdbmap,

		_struct:    _struct,
		_pstruct:   _pstruct,
		_nilstruct: nil,

		_chan:  make(chan int),
		_pchan: &_pchan,
		_func:  _func,

		_interface: nil,

		_nest: &AA{123456, s_NetInfo{"net", 100}, &S_DBInfo{"db", 200}},
	}

	fmt.Println(SprintStruct(cfg, ">>", "    "))

	fstest.PrintTestEnd()
}
