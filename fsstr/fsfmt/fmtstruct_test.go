package fsfmt

import "fmt"
import "testing"
import "fsky.pro/fstest"

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
	NetInfo s_NetInfo
	DBInfo  S_DBInfo
	LogPath string
	AA
}

func TestSprintStruct(t *testing.T) {
	fstest.PrintTestBegin("SprintStruct")

	c := Config{
		NetInfo: s_NetInfo{"172.16.146.124", 100},
		DBInfo:  S_DBInfo{"aa bb cc", 3000},
		LogPath: "xxxxx xxxxxx",
	}
	fmt.Println(SprintStruct(c, ">> ", "--"))

	fstest.PrintTestEnd()
}
