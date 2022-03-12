/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: sql info build return
@author: fanky
@version: 1.0
@date: 2022-01-04
**/

package fssql

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// -------------------------------------------------------------------
// SQLInfo
// -------------------------------------------------------------------
type S_SQLInfo struct {
	err      error
	sqlText  string
	InValues []interface{}
}

var _sqlInfoType = reflect.TypeOf(S_SQLInfo{})

func newSQLInfo(err error, sqlText string) *S_SQLInfo {
	return &S_SQLInfo{
		err:     err,
		sqlText: sqlText,
	}
}

func (this *S_SQLInfo) Type() string {
	return fmt.Sprintf("%s", _sqlInfoType)
}

func (this *S_SQLInfo) Err() error {
	return this.err
}

func (this *S_SQLInfo) SQLText() string {
	return this.sqlText
}

func (this *S_SQLInfo) FmtSQLText(indent string) string {
	re := regexp.MustCompile(`\?`)
	count := len(this.InValues)
	if count == 0 {
		return indent + this.sqlText
	}

	inValues := []string{}
	index := 0
	sqltx := re.ReplaceAllStringFunc(this.sqlText, func(e string) string {
		index++
		if index <= count {
			inValues = append(inValues, fmt.Sprintf("%d:%#v", index, this.InValues[index-1]))
		} else {
			inValues = append(inValues, "<???>")
		}
		return fmt.Sprintf("?[%d]", index)
	})
	strInValues := "{" + strings.Join(inValues, ", ") + "}"
	return fmt.Sprintf("%s%s;\n%s  ?: %s", indent, sqltx, indent, strInValues)
}

// -------------------------------------------------------------------
// Select SQLInfo
// -------------------------------------------------------------------
type S_SelectInfo struct {
	S_SQLInfo
	selTable   *S_Table
	selMembers []*S_Member
}

func newSelectInfo(sqlInfo *S_SQLInfo, selTable *S_Table, members []*S_Member) *S_SelectInfo {
	return &S_SelectInfo{
		S_SQLInfo:  *sqlInfo,
		selTable:   selTable,
		selMembers: members,
	}
}

// 创建传出对象
func (this *S_SelectInfo) CreateOutObject() (obj interface{}, mptrs []interface{}, err error) {
	if this.selTable == nil || len(this.selMembers) == 0 {
		err = fmt.Errorf("not select objects in this building sql")
		return
	}
	mptrs = []interface{}{}
	pobj := reflect.New(this.selTable.tobj)
	obj = pobj.Interface()
	vobj := pobj.Elem()
	for _, m := range this.selMembers {
		mptrs = append(mptrs, m.valuePtr(vobj))
	}
	return
}

// 获取指定对象的传出成员指针
func (this *S_SelectInfo) ScanMembers(obj interface{}) ([]interface{}, error) {
	if this.selTable == nil || len(this.selMembers) == 0 {
		return nil, fmt.Errorf("this SelectInfo is not support output object")
	}
	tobj := reflect.TypeOf(obj)
	if tobj == nil {
		return nil, fmt.Errorf("output object is not allow be a nil value")
	}
	if tobj.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("the output object must be a object pointer of type %v", this.selTable.tobj)
	}
	vobj := reflect.ValueOf(obj)
	if vobj.IsNil() {
		return nil, fmt.Errorf("output object is not allow be a nil value")
	}
	vobj = vobj.Elem()
	outMembers := []interface{}{}
	for _, m := range this.selMembers {
		outMembers = append(outMembers, m.valuePtr(vobj))
	}
	return outMembers, nil
}

// -------------------------------------------------------------------
// Insert SQLInfo
// -------------------------------------------------------------------
type S_ExecInfo struct {
	S_SQLInfo
}

func newExecInfo(sqlInfo *S_SQLInfo) *S_ExecInfo {
	return &S_ExecInfo{
		S_SQLInfo: *sqlInfo,
	}
}
