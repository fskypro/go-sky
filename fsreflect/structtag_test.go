package fsreflect

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

type AA struct {
	value string `db:"xx"`
}

type BB struct {
	value int `db:"yy"`
}

type CC struct {
	AA
	*BB

	OpenID  string `db:"openid"`
	UnionID string `db:"unionid"`

	Name     string `db:"name"`
	Head     string `db:"head"`
	Sex      int8   `db:"sex"`
	Province string `db:"province"`
	City     string `db:"city"`

	Score     uint32 `db:"score"`
	Grade     uint32 `db:"grade"`
	LoginTime string `db:logintime`
}

func TestFieldTagMap(t *testing.T) {
	fstest.PrintTestBegin("FieldTagMap")

	c := CC{}
	fmt.Println(FieldTagsMap(c, "db"))
	fmt.Println(TagFieldsMap(&c, "db"))

	fstest.PrintTestEnd()
}

func TestTagFieldMap(t *testing.T) {
	fstest.PrintTestBegin("TagFieldMap")

	c := CC{}
	fmt.Println(FieldTagsMap(c, "db"))
	fmt.Println(TagFieldsMap(&c, "db"))

	fstest.PrintTestEnd()
}
