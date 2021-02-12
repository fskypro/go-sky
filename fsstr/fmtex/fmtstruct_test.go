package fmtex

import (
	"fmt"
	"strings"
	"testing"

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
}

type Config struct {
	i        int
	pi       *int
	NetInfo  s_NetInfo
	DBInfo   *S_DBInfo
	LogPath  string
	PLogPath *string
	AA
}

func TestSprintStruct(t *testing.T) {
	fstest.PrintTestBegin("SprintStruct")
	ii := 200
	logPath := "yyyyy, yyyyy"
	c := Config{
		i:        100,
		pi:       &ii,
		NetInfo:  s_NetInfo{"172.16.146.124", 100},
		DBInfo:   &S_DBInfo{"aa bb cc", 3000},
		LogPath:  "xxxxx, xxxxxx",
		PLogPath: &logPath,
	}
	fmt.Println(SprintStruct(c, ">> ", "--"))
	fmt.Println(strings.Repeat("-", 20))
	fmt.Println(SprintStruct(&c, ">> ", "--"))
	fmt.Println(strings.Repeat("-", 20))

	fstest.PrintTestEnd()
}
