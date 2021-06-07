package fsreflect

import (
	"fmt"
	"testing"

	"fsky.pro/fsreflect/test"
	"fsky.pro/fsstr/fmtex"
	"fsky.pro/fstest"
)

func TestGetFieldValue(t *testing.T) {
	fstest.PrintTestBegin("GetFieldValue")

	obj := test.NewObj()
	Aa, _ := GetFieldValue(obj, "Aa")
	bb, _ := GetFieldValue(obj, "bb")
	Pcc, _ := GetFieldValue(obj, "Pcc")
	pcc, _ := GetFieldValue(obj, "pcc")
	DD, _ := GetFieldValue(obj, "DD")
	dd, _ := GetFieldValue(obj, "dd")
	pdd, _ := GetFieldValue(obj, "pdd")
	nildd, _ := GetFieldValue(obj, "nildd")
	fmt.Println("obj.Aa = ", Aa.(string))
	fmt.Println("obj.bb = ", bb.(string))
	fmt.Println("obj.Pcc = ", *(Pcc.(*string)))
	fmt.Println("obj.pcc = ", *(pcc.(*string)))
	fmt.Println("obj.DD = ", DD.(test.DD))
	fmt.Println("obj.dd = ", dd.(test.DD))
	fmt.Println("obj.pdd = ", *(pdd.(*test.DD)))
	fmt.Println("obj.nildd = ", nildd)

	fstest.PrintTestEnd()
}

func TestSetFieldValue(t *testing.T) {
	fstest.PrintTestBegin("GetFieldValue")

	obj := test.NewObj()
	Pcc := "cccc"
	pcc := "dddd"
	nildd := &test.DD{3000}
	SetFieldValue(obj, "Aa", "aaaa")
	SetFieldValue(obj, "bb", "bbbb")
	SetFieldValue(obj, "Pcc", &Pcc)
	SetFieldValue(obj, "pcc", &pcc)
	SetFieldValue(obj, "DD", test.DD{1000})
	SetFieldValue(obj, "dd", test.DD{2000})
	SetFieldValue(obj, "pdd", nil)
	SetFieldValue(obj, "nildd", nildd)

	fmt.Println("obj.Aa = ", obj.Aa)
	fmt.Println("obj.bb = ", obj.GetBB())
	fmt.Println("obj.Pcc = ", *obj.Pcc)
	fmt.Println("obj.pcc = ", *obj.GetPCC())
	fmt.Println("obj.DD = ", obj.DD)
	fmt.Println("obj.dd = ", obj.GetDD())
	fmt.Println("obj.pdd = ", obj.GetPDD())
	fmt.Println("obj.nildd = ", *obj.GetNilDD())

	fstest.PrintTestEnd()
}

func TestCopyStructObject(t *testing.T) {
	fstest.PrintTestBegin("CopyStructObject")

	type A struct {
		V string
	}

	type S struct {
		V1 string
		v2 int
		v3 *A
	}

	s := &S{
		V1: "xxx",
		v2: 100,
		v3: &A{"yyy"},
	}
	var ss *S = new(S)

	err := CopyStructObject(ss, s)
	if err != nil {
		fmt.Println("copy struct error:", err.Error())
		return
	}
	fmt.Println(fmtex.SprintStruct(ss, nil))

	s.v3.V = "zzz"
	fmt.Println(fmtex.SprintStruct(ss, nil))

	fstest.PrintTestEnd()
}
