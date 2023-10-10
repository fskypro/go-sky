package fssql

import (
	"fmt"
	"testing"
)

type S_Test struct {
	A string `db:"a"`
	B int    `db:"b"`
	C int    `db:"c"`
}

var tb *S_Table

func init() {
	tb, _ = NewTable("table", new(S_Test))
}

func TestRe(t *testing.T) {
	test := new(S_Test)
	sql := LeadSQL('%').SQL(`select %TO{A,C} from %[2]TN where %[2]TM{A}=%[3]q `, test, tb, "abc").
		SQL("AND #TM{C}=#v", tb, 200)
	fmt.Println(sql.Error)
	fmt.Println(11111, sql.SQLTxt)
	fmt.Println(22222, sql.Inputs)

	sql = SQL(`update %TN set %TE-.{A,B}`, tb, test)
	fmt.Println(3333, sql.SQLTxt)

	sql = SQL(
		"WITH RECURSIVE sub_nodes AS ("+
			"SELECT %[1]TM{A, B} FROM %[1]TN WHERE %[1]TM{A} = %[2]v "+
			"UNION "+
			"SELECT %[1]TM(t){A,A} FROM %[1]TN t INNER JOIN sub_nodes ON %[1]TM(t){B} = %[1]TM(sub_nodes){A} "+
			") "+
			"DELETE FROM %[1]TN WHERE id IN (SELECT id FROM sub_nodes)", tb, 100)
	fmt.Println(44444, sql.SQLTxt)
}
