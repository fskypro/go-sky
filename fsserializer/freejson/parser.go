/**
@copyright: fantasysky 2016
@brief: 解释 json 字符串
@author: fanky
@version: 1.0
@date: 2019-05-31
**/

package freejson

//import "fmt"
import (
	"errors"
	"strconv"
	"strings"
)

// json 转移符对应 go 的转义符
var _byte2trans = map[byte]byte{
	'\\': '\\', // 反斜杠
	'b':  '\b', // 退格符
	'r':  '\r', // 回车
	'n':  '\n', // 换行
	'f':  '\f', // 换页
}

type s_Parser struct {
	data  []byte
	count int
	index int

	subParsers []func(*s_Parser) (I_Value, error)
}

func newParser(jbytes []byte) *s_Parser {
	parser := &s_Parser{
		data:  jbytes,
		count: len(jbytes),
	}
	parser.subParsers = []func(*s_Parser) (I_Value, error){
		_parseString,
		_parseNumber,
		_parseBool,
		_parseNull,
		_parseList,
		_parseObject,
	}
	return parser
}

// -------------------------------------------------------------------
// module private
// -------------------------------------------------------------------
func (this *s_Parser) _newParseError() *ParseError {
	var row, col int = 0, 0
	isr := false
	for index := 0; index <= this.index; index++ {
		char := this.data[index]

		// 前一个字符是 \r
		if isr {
			if char == '\n' {
				continue
			} else if char != '\r' {
				isr = false
			}
		}

		if char == '\r' {
			isr = true
		}

		if char == '\r' || char == '\n' {
			row += 1
			col = 0
		} else {
			col += 1
		}
	}
	return newParseError(this.index, row, col)
}

// --------------------------------------------------
// 是否解释文档结束
func (this *s_Parser) _isParseEnd() bool {
	return this.index >= this.count
}

// 读取当前游标后 n 个字符，文档结束则返回 nil
func (this *s_Parser) _look(n int) []byte {
	if this.index >= this.count {
		return nil
	}
	endIndex := this.index + n
	if endIndex > this.count {
		endIndex = this.count
	}
	return this.data[this.index:endIndex]
}

// --------------------------------------------------
// 游标下跳 1 个字符，并返回该字符，如后面没有字符（文档结束），则返回 0
func (this *s_Parser) _pickOne() byte {
	if this.index < this.count {
		char := this.data[this.index]
		this.index += 1
		return char
	}
	return 0
}

// 移动游标到指定的一组字符中，任意一个字符后，并返回中间跨厉的字符串（不包含指定的字符）
// skipEnd 表示指针是否跳过结束符
func (this *s_Parser) _pickUntil(skipEnd bool, chars ...byte) []byte {
	start := this.index
	for this.index < this.count {
		char := this.data[this.index]
		for _, c := range chars {
			if c == char {
				str := this.data[start:this.index]
				if skipEnd {
					this.index += 1
				}
				return str
			}
		}
		this.index += 1
	}
	return nil
}

// 提取字符串
// 以下是转移符：
// \\ 反斜杠
// \b 退格
// \r 回车
// \n 换行
// \f 换页
// \" 双引号
func (this *s_Parser) _pickQuoteString() ([]byte, error) {
	chars := this._look(1)
	if chars == nil || len(chars) != 1 {
		return nil, nil
	}
	quote := chars[0]
	if quote != '"' && quote != '`' {
		return nil, nil
	}
	this._skip(1)

	str := make([]byte, 0)
	var char byte
	slash := false // 上一个字符是否是反斜杠
	for this.index < this.count {
		char = this.data[this.index]
		this.index += 1

		// 进入转义
		if slash {
			c, ok := _byte2trans[char]
			if ok {
				str = append(str, c)
			} else if char == quote {
				str = append(str, quote)
			} else {
				str = append(str, '\\')
				str = append(str, char)
			}
			slash = false
			continue
		}

		// 字符串结束
		if char == quote {
			return str, nil
		} else if char == '\\' {
			slash = true
		} else {
			str = append(str, char)
		}
	}
	return nil, errors.New("")
}

// --------------------------------------------------
// 游标跳过指定数量
func (this *s_Parser) _skip(n int) {
	this.index += n
	if this.index > this.count {
		this.index = this.count
	}
}

// 游标跳过当前行(跳到下一行)
func (this *s_Parser) _skipLine() {
	for this.index < this.count {
		char := this.data[this.index]
		this.index += 1
		if char == '\r' {
			if char == '\n' {
				this.index += 1
			}
			break
		} else if char == '\n' {
			break
		}
	}
}

// 游标跳过所有空白字符
func (this *s_Parser) _skipSpace() {
	for this.index < this.count {
		char := this.data[this.index]
		if char == ' ' || char == '\t' ||
			char == '\r' || char == '\n' {
			this.index += 1
		} else {
			break
		}
	}
}

// 跳过注释
func (this *s_Parser) _skipComment() bool {
	chars := this._look(2)
	if chars == nil || len(chars) != 2 || chars[0] != '/' {
		return false
	}

	if chars[1] == '/' { // 单行注释
		this._skipLine() // 忽略整行
		return true
	}
	if chars[1] != '*' { // 不是多行注释
		return true
	}

	// 寻找多行注释结束符
	this._skip(2)
	for this.index < this.count {
		chars = this._look(2)
		if chars[0] == '*' && chars[1] == '/' {
			this._skip(2)
			return true
		}
		this.index += 1
	}
	return true
}

// 跳过无效内容，包括空白字符和注释
func (this *s_Parser) _skipInvalids() {
	for {
		this._skipSpace()
		if !this._skipComment() {
			break
		}
	}
}

// ------------------------------------------------------------
// 解释键
func (this *s_Parser) _parseKey() []byte {
	key, _ := this._pickQuoteString()
	if key == nil {
		return nil
	}

	this._skipInvalids()
	if this._pickOne() != ':' {
		return nil
	}
	return key
}

// 解释字符串值
func _parseString(this *s_Parser) (value I_Value, err error) {
	str, err := this._pickQuoteString()
	if str == nil {
		return nil, err
	}
	return newString(str), nil
}

// 解释数学数值
// int64：- 或 + 或 数字开头
// uint64：u+数字(十进制)；b+[0|1](二进制)；0+小于7的数字(八进制)；0x|0X+[0-9|a-b](十六进制)
// float64：+ 或 - + 带小数点的数字
func _parseNumber(this *s_Parser) (I_Value, error) {
	prefix := string(this._look(2))
	if len(prefix) != 2 {
		return nil, nil
	}
	first := prefix[0]

	// 无符号整数
	if first == 'u' { // 十进制
		this._skip(1)
		strValue := string(this._pickUntil(false, ',', ' ', '\t', '\r', '\n', '}', ']'))
		value, err := strconv.ParseUint(strValue, 10, 64)
		if err != nil {
			return nil, err
		}
		return NewUInt64(value), nil
	}
	if prefix == "0b" { // 二进制
		this._skip(2)
		strValue := string(this._pickUntil(false, ',', ' ', '\t', '\r', '\n', '}', ']'))
		value, err := strconv.ParseUint(strValue, 2, 64)
		if err != nil {
			return nil, err
		}
		return NewUInt64(value), nil
	}
	if prefix == "0X" || prefix == "0x" { // 十六进制
		this._skip(2)
		bsValue := this._pickUntil(false, ',', ' ', '\t', '\r', '\n', '}', ']')
		if bsValue == nil {
			return nil, errors.New("")
		}
		value, err := strconv.ParseUint(string(bsValue), 16, 64)
		if err != nil {
			return nil, err
		}
		return NewUInt64(value), nil
	}
	if first == '0' { // 八进制或小于 1 的浮点数
		strValue := string(this._pickUntil(false, ',', ' ', '\t', '\r', '\n', '}', ']'))
		if strings.Index(strValue, ".") >= 0 { // 小于 1 的浮点数
			value, err := strconv.ParseFloat(strValue, 64)
			if err != nil {
				return nil, err
			}
			return NewFloat64(value), nil
		} else { // 八进制
			value, err := strconv.ParseUint(strValue, 8, 64)
			if err != nil {
				return nil, err
			}
			return NewUInt64(value), nil
		}
	}

	// 有符号整数或者浮点数
	if first == '+' || first == '-' || ('0' < first && first <= '9') || first == '.' {
		strValue := string(this._pickUntil(false, ',', ' ', '\t', '\r', '\n', '}', ']'))
		if strings.Index(strValue, ".") >= 0 { // 浮点数
			value, err := strconv.ParseFloat(strValue, 64)
			if err != nil {
				return nil, err
			}
			return NewFloat64(value), nil
		} else { // 有符号整数
			value, err := strconv.ParseInt(strValue, 10, 64)
			if err != nil {
				return nil, err
			}
			return NewInt64(value), nil
		}
	}

	return nil, nil
}

// 解释布尔值
func _parseBool(this *s_Parser) (I_Value, error) {
	if string(this._look(4)) == "true" {
		this._skip(4)
		return NewBool(true), nil
	}
	if string(this._look(5)) == "false" {
		this._skip(5)
		return NewBool(false), nil
	}
	return nil, nil
}

// 解释 null 值
func _parseNull(this *s_Parser) (I_Value, error) {
	if string(this._look(4)) == "null" {
		this._skip(4)
		return NewNull(), nil
	}
	return nil, nil
}

// 解释列表
func _parseList(this *s_Parser) (I_Value, error) {
	chars := this._look(1)
	if chars[0] != '[' {
		return nil, nil
	} else {
		this._skip(1)
	}

	list := NewList()
	for {
		this._skipInvalids()
		if this._isParseEnd() {
			return nil, errors.New("")
		}

		chars = this._look(1)
		if chars == nil {
			return nil, errors.New("")
		}

		if chars[0] == ']' { // 列表结束
			this._skip(1)
			return list, nil
		}

		var err error
		var value I_Value = nil
		for _, parser := range this.subParsers {
			value, err = parser(this)
			if err != nil {
				return nil, err
			}
			if value != nil {
				list.Add(value)
				break
			}
		}
		if value == nil {
			return nil, errors.New("")
		}

		this._skipInvalids()
		char := this._pickOne()
		if char == ']' {
			return list, nil
		} else if char != ',' {
			return nil, errors.New("")
		}
	}
}

// 解释子节点
func (this *s_Parser) _parseItem() (key []byte, value I_Value) {
	key = this._parseKey()
	if key == nil {
		return
	}

	// 跳过冒号后的注释和空白
	this._skipInvalids()
	if this._isParseEnd() {
		return
	}

	// 解释值
	var err error
	for _, parser := range this.subParsers {
		value, err = parser(this)
		if err != nil { // 解释有错
			return
		}
		if value != nil {
			break
		}
	}
	return
}

// 解释对象
func _parseObject(this *s_Parser) (I_Value, error) {
	chars := this._look(1)
	if chars[0] != '{' {
		return nil, nil
	} else {
		this._skip(1)
	}

	obj := NewObject()
	for this.index < this.count {
		// 去掉注释和空白
		this._skipInvalids()

		chars := this._look(1)
		// 文档结束
		if chars == nil {
			return nil, errors.New("")
		}

		// object 结束
		if chars[0] == '}' {
			this._skip(1)
			return obj, nil
		}

		// 解释子节点
		key, value := this._parseItem()
		if key == nil || value == nil {
			return nil, errors.New("")
		}
		obj.Add(string(key), value)

		// 检查子节点结束符号
		this._skipInvalids()
		char := this._pickOne()
		if char == '}' {
			return obj, nil
		} else if char != ',' {
			return nil, errors.New("")
		}
	}
	return nil, errors.New("")
}

// -------------------------------------------------------------------
// package public
// -------------------------------------------------------------------
func (this *s_Parser) parse() (value I_Value, err error) {
	// 去掉前面的注释和空白字符
	this._skipInvalids()
	if this._isParseEnd() {
		err = this._newParseError()
		return
	}

	value = nil
	for _, parser := range this.subParsers {
		value, err = parser(this)
		if err != nil {
			continue
		}
		if value != nil {
			break
		}
	}
	if value == nil {
		return nil, this._newParseError()
	}

	// 创建根节点
	//obj, err := _parseObject(this)
	//if err != nil || obj == nil {
	//err = this._newParseError()
	//return
	//}

	// 检查文档末尾
	this._skipInvalids()
	if !this._isParseEnd() {
		err = this._newParseError()
	}
	return
}
