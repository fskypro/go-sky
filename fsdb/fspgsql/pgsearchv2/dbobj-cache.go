/**
@copyright: yxht 2015
@website: http://www.exht.com
@brief: db object cache
@author: fanky
@version: 1.0
@date: 2025-07-18
**/

package pgsearchv2

import (
	"fmt"
	"reflect"
	"sync"
	"unicode"

	"fsky.pro/fspgsql/fssql"
	"fsky.pro/fsreflect"
)

// -------------------------------------------------------------------
// 数据库对象名称缓存
// -------------------------------------------------------------------
type s_DBObjCache struct {
	mu            sync.RWMutex
	jkeyMembers   map[reflect.Type]map[string]*fssql.S_Member // {<对象结构类>: {<json-key>: <*fssql.S_Member>}}
	cndParserPkgs sync.Map                                    // {<对象结构类>: true}
}

var dbObjCache = &s_DBObjCache{
	jkeyMembers: map[reflect.Type]map[string]*fssql.S_Member{},
}

func (this *s_DBObjCache) getFieldMember(table *fssql.S_Table, dbObj any, jkey string) *fssql.S_Member {
	t := reflect.TypeOf(dbObj)
	if t == nil {
		return nil
	}
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	mb := func() *fssql.S_Member {
		this.mu.RLock()
		defer this.mu.RUnlock()
		if this.jkeyMembers[t] == nil {
			return nil
		}
		return this.jkeyMembers[t][jkey]
	}()

	if mb != nil {
		// 缓存中能找到
		return mb
	}

	// 遍历对象成员
	// 根据查找对象中导出的 json 字段，获取对应的数据库列字段
	fsreflect.TrivalStructMembers(dbObj, true, func(info *fsreflect.S_TrivalStructInfo) bool {
		jsonTag := info.Field.Tag.Get("json")
		if jsonTag == "-" {
			return true
		} else if jsonTag == "" {
			jsonTag = info.Field.Name
		}
		if jsonTag == "" {
			return true
		}
		if unicode.IsLower(rune(info.Field.Name[0])) {
			return true
		}
		if jsonTag != jkey {
			return true
		}
		mb = table.Member(info.Field.Name)
		return false
	})
	if mb != nil {
		this.mu.Lock()
		defer this.mu.Unlock()
		if this.jkeyMembers[t] == nil {
			this.jkeyMembers[t] = map[string]*fssql.S_Member{jkey: mb}
		} else {
			this.jkeyMembers[t][jkey] = mb
		}
	}
	return mb
}

// 获取数据库对象的解释搜索参数的成员方法
// 如果找到，则将搜索参数的解释工作交给成员方法
func (this *s_DBObjCache) getCndParser(dbObj any, cndInfo *S_PGCndInfo) (m reflect.Value, err error) {
	const parserName = "ParseSearchCnd"

	val := reflect.ValueOf(dbObj)
	method := val.MethodByName(parserName)
	if !method.IsValid() {
		// 对象没有实现解释方法
		return
	}

	t := reflect.TypeOf(dbObj)
	if t == nil {
		return
	}
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 判断是否已经校验过参数和返回值的合法性
	_, valid := this.cndParserPkgs.Load(t)
	if !valid {
		methodType := method.Type()

		// 检查签名：输入数量 == 1，输出数量 == 2
		if methodType.NumIn() != 2 || methodType.NumOut() != 3 {
			goto L
		}

		// 检查参数 1 的类型是 string
		argType := methodType.In(0)
		expectedArg := reflect.TypeOf("")
		if argType != expectedArg {
			goto L
		}

		// 检查参数 2 的类型是 *dbsearchv2.S_CndInfo
		argType = methodType.In(1)
		expectedArg = reflect.TypeOf(cndInfo)
		if argType != expectedArg {
			goto L
		}

		// 检查返回类型 ([]*dbsearchv2.S_CndInfo, error)
		if methodType.Out(0) != reflect.TypeOf("") ||
			methodType.Out(1) != reflect.TypeOf([]any{}) ||
			methodType.Out(2) != reflect.TypeOf((*error)(nil)).Elem() {
			goto L
		}
		this.cndParserPkgs.Store(t, true)
		valid = true
	}

L:
	if !valid {
		err = fmt.Errorf("method %q of %q must be like as: %s(string, *pgsearchv2.S_PGCndInfo) (string, []any, error)",
			parserName, reflect.TypeOf(dbObj), parserName)
		return
	}
	m = method
	return
}
