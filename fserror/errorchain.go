/**
@copyright: fantasysky 2016
@brief: 实现错误缓冲链
@author: fanky
@version: 1.0
@date: 2019-01-02
**/

package fserror

import "sync"
import "strings"
import "fsky.pro/fsenv"

// S_ErrorChain 错误链
type S_ErrorChain struct {
	sync.Mutex
	Errors []error
}

func NewErrorChain() *S_ErrorChain {
	return &S_ErrorChain{
		Errors: []error{},
	}
}

// 缓存的错误个数
func (this *S_ErrorChain) Count() int {
	return len(this.Errors)
}

// Error 返回错误链条中的所有错误，并以换行符分隔
func (this *S_ErrorChain) Error() string {
	sb := new(strings.Builder)
	for _, err := range this.Errors {
		sb.WriteString(err.Error())
		sb.WriteString(fsenv.Endline)
	}
	return sb.String()
}

// Append 追加一个错误
func (this *S_ErrorChain) Append(err error) {
	this.Lock()
	defer this.Unlock()
	this.Errors = append(this.Errors, err)
}

// AppendStrError 追加一个可字符串格式化的错误
func (this *S_ErrorChain) AppendStrErrorf(format string, args ...interface{}) {
	err := StrErrorf(format, args...)
	this.Lock()
	defer this.Unlock()
	this.Errors = append(this.Errors, err)
}
