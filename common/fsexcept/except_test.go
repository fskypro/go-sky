package fsexcept

import (
	"fmt"
	"testing"
)

func dev(v int) int {
	if v == 0 {
		Throw("zero devid")
	}
	return 100 / v
}

func test(v int) {
	var value int
	S_Except{
		Try: func(S_Except) {
			value = dev(v)
			fmt.Println("result:", value)
		},
		Catch: func(e any) {
			fmt.Printf("except: %v\n", e)
		},
		Finally: func() {
			fmt.Printf("finally!\n\n")
		},
	}.Do()
}

func Test(t *testing.T) {
	test(1)
	test(0)
}
