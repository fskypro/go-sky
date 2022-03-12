package fsstr

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	s := []interface{}{"xxxx", "yyyy", nil, 12}
	str := JoinFunc(s, ", ", func(a interface{}) string {
		return fmt.Sprintf("%#v", a)
	})
	fmt.Println(str)
}
