package freejson

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

var s = `
  // 这是单行注释
  /*
  这是多行注释 
  */

  {
      "a1": {
          "b1": null ,
          "b1": "aaaaaa",
          "c1": ` + "`" + `aa"aaaa` + "`" + `     // 也可以用反引号作为字符串的括弧，这样，字符串中的双引号不需要斜杠作为转义
      },

      // 以下写法都是无符号整数(go 中类型都是 uint64)
      "uint_10": u213,        // 10 进制          
      "uint_16": 0x344,       // 16 进制
      "uint_2" : 0b10101,     // 2 进制
      "uint_8" : 012317,      // 8 进制(0 开头)

      // 以下是有符号整数(go 中类型都是 int64)
      "int_64" : 123456,      // 10 进制
      "-int_64": -100,        // 10 进制负数
      /* 注意：
          有符号整数只有 10 进制的表示法
      */

      // 以下是浮点数（go 中都是 float64）
      "float": 0.453, 
      "-float": -2.32,

      // 列表
      "list": ["123", 546, {"k1": "字符串", "k2": -3425245}, /* 随意嵌入多行注释 */],

      // 嵌套多层
      "k": {
          "kk1": {
              "kkk1": {
                  "kkkk1": {
                      "随意": ["e1", 123, "随意"]
                  }
              },
              "kkk2": "很好"
          }
      },
      // 最后一个元素后面，要不要 , 号都无所谓
  }
`

func TestFromString(t *testing.T) {
	fstest.PrintTestBegin("FromString")

	// 解释
	v, err := FromString(s)
	r := v.(*S_Object)
	if err != nil {
		fmt.Println("parse json string fail: ", err.Error())
		return
	}

	// 打印
	r.For(func(k string, v I_Value) bool {
		fmt.Printf("%s: %v\n", k, v)
		return true
	})

	// 保存成文件
	fmtInfo := NewFmtInfo()
	//fmtInfo.IndentList = true // 是否对列表每个元素重起一行格式化
	fmtInfo.Indent = "\t"
	err = Save(r, "./test.json", fmtInfo)
	if err != nil {
		fmt.Printf("\nsave json file fail: %s\n", err.Error())
		return
	} else {
		fmt.Println("\nsave json file success!")
	}

	// 重新读取文件
	vv, err := Load("./test.json")
	r = vv.(*S_Object)
	if err != nil {
		fmt.Printf("\n load json file fail: %s\n", err.Error())
		return
	}

	// 打印
	r.For(func(k string, v I_Value) bool {
		fmt.Printf("%s: %v\n", k, v)
		return true
	})

	fstest.PrintTestEnd()
}

func TestList(t *testing.T) {
	fstest.PrintTestBegin("FromList")
	doc := `["aa", "bb", 123, 567, {"k1": 345, "k2": 678}]`
	v, err := FromString(doc)
	if err != nil {
		fmt.Println("load list fail: ", err.Error())
		return
	}

	l := v.(*S_List)
	fmt.Println(l.String())

	l.For(func(index int, elem I_Value) bool {
		fmt.Printf("%v: %v\n", index, elem)
		return true
	})

	err = Save(v, "./list.json", nil)
	if err != nil {
		fmt.Println("save list fail: ", err.Error())
	}

	fstest.PrintTestEnd()
}
