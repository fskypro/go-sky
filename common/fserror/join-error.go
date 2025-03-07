/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: error joiner
@author: fanky
@version: 1.0
@date: 2024-12-11
**/

package fserror

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"unsafe"
)

type s_FLError struct {
	prefix string
	err    error
}

func newFLError(deep int, cutRoot string, format string, args ...any) error {
	if deep < 1 {
		return fmt.Errorf(format, args...)
	}
	_, file, line, ok := runtime.Caller(deep + 1)
	if !ok {
		return fmt.Errorf(format, args...)
	}

	if cutRoot != "" {
		prefix, err := filepath.Rel(cutRoot, file)
		if err != nil {
			prefix = file
		}
		prefix = strings.TrimSuffix(prefix, ".go")
		prefix = fmt.Sprintf("%s:%d: ", prefix, line)
		return s_FLError{
			prefix: prefix,
			err:    fmt.Errorf(format, args...),
		}
	}
	return s_FLError{
		prefix: fmt.Sprintf("%s:%d: ", file, line),
		err:    fmt.Errorf(format, args...),
	}
}

func (self s_FLError) errorWithFileLine() string {
	return self.prefix + self.err.Error()
}

func (self s_FLError) Error() string {
	return self.err.Error()
}

// -------------------------------------------------------------------
// JoinError
// -------------------------------------------------------------------
type S_JError struct {
	indent string
	errs   []error
}

func (this *S_JError) Join(errs ...error) error {
	n := 0
	for _, err := range errs {
		if err != nil {
			n++
		}
	}
	if n == 0 {
		return this
	}
	for _, err := range errs {
		if err == nil {
			continue
		}
		switch err.(type) {
		case *S_JError:
			this.errs = append(this.errs, err.(*S_JError).errs...)
		default:
			this.errs = append(this.errs, errs...)
		}
	}
	return this
}

func (this *S_JError) Error() string {
	b := []byte{}
	for i, err := range this.errs {
		b = append(b, '\n')
		b = append(b, []byte(strings.Repeat(this.indent, i+1))...)
		switch err.(type) {
		case s_FLError:
			b = append(b, err.(s_FLError).errorWithFileLine()...)
		default:
			b = append(b, err.Error()...)
		}
	}
	return unsafe.String(&b[0], len(b))
}

func (this *S_JError) Unwrap() []error {
	return this.errs
}

func JErrorf(format string, args ...any) *S_JError {
	return &S_JError{
		indent: "",
		errs:   []error{fmt.Errorf(format, args...)},
	}
}

// 带逐层缩进
func JIndentErrorf(indent string, format string, args ...any) *S_JError {
	return &S_JError{
		indent: indent,
		errs:   []error{fmt.Errorf(format, args...)},
	}
}

// 带文件名和行号
func JFLErrorf(indent string, deep int, format string, args ...any) *S_JError {
	return &S_JError{
		indent: indent,
		errs:   []error{newFLError(deep+1, "", format, args...)},
	}
}

// 带去掉根部的文件名和行号
func JCFLErrorf(indent string, deep int, cutRoot string, format string, args ...any) *S_JError {
	return &S_JError{
		indent: indent,
		errs:   []error{newFLError(deep+1, cutRoot, format, args...)},
	}
}
