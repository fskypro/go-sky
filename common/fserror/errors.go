/**
@copyright: fantasysky 2016
@brief: errors set
@author: fanky
@version: 1.0
@date: 2020-05-02
**/

package fserror

import "fmt"
import "fsky.pro/fsenv"

type S_Errors struct {
	errs []error
}

func NewErrors() *S_Errors {
	return &S_Errors{
		errs: []error{},
	}
}

func (this *S_Errors) Add(err error) {
	this.errs = append(this.errs, err)
}

func (this *S_Errors) Addf(msg string, args ...interface{}) {
	this.errs = append(this.errs, fmt.Errorf(msg, args...))
}

func (this *S_Errors) Errors() []error {
	return this.errs
}

func (this *S_Errors) Count() int {
	return len(this.errs)
}

func (this *S_Errors) HasError(err error) bool {
	for _, e := range this.errs {
		if e == err {
			return true
		}
	}
	return false
}

func (this *S_Errors) FmtErrors() string {
	if len(this.errs) == 0 {
		return ""
	}

	var msg string
	first := this.errs[0]
	if e, ok := first.(I_Error); ok {
		msg = "+ " + e.Error()
	} else {
		msg = "+ " + first.Error()
	}

	for _, err := range this.errs[1:] {
		if e, ok := err.(I_Error); ok {
			msg = msg + fsenv.Endline + "+ " + e.Error()
		} else {
			msg = msg + fsenv.Endline + "+ " + err.Error()
		}
	}
	return msg
}
