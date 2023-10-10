package fsalgorithm

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

type s_Obj struct {
	Value int
}

func TestSort(t *testing.T) {
	fstest.PrintTestBegin("SortFunc")
	defer fstest.PrintTestEnd()

	l := []s_Obj{
		s_Obj{4},
		s_Obj{0},
		s_Obj{2},
		s_Obj{5},
		s_Obj{3},
	}

	SortFunc(l[1:], func(s1, s2 s_Obj) bool { return s1.Value < s2.Value })

	fmt.Println(l)
}
