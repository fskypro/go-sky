package fserror

import "fmt"
import "errors"
import "testing"
import "fsky.pro/fstest"

type s_BaseError struct {
	msg string
}

func (this *s_BaseError) Error() string {
	return this.msg
}

var baseError = new(s_BaseError)

func newError1() error {
	return Wrapf(baseError, "new error1")
}

type s_Error2 struct {
	S_Error
}

func newError2() I_Error {
	return &s_Error2{
		S_Error: *Wrapf(newError1(), "new error2"),
	}
}

func newError3() I_Error {
	return Wrapf(newError2(), "new error3")
}

func TestError(t *testing.T) {
	fstest.PrintTestBegin("fserror.Error")
	err := newError3()
	fmt.Println("errors.Is(newError2(), new(s_BaseError)) = ", errors.Is(err, new(s_BaseError)))
	fmt.Println("errors.Is(newError2(), baseError) = ", errors.Is(err, baseError))

	var tmp *s_BaseError
	fmt.Println("errors.As(newError2(), *s_BaseError) = ", errors.As(err, &tmp))
	fmt.Println("error message:")
	fmt.Println(err.Error())
	fstest.PrintTestEnd()
}

func TestErrors(t *testing.T) {
	fstest.PrintTestBegin("fserror.Errors")
	errs := NewErrors()
	errs.Addf("the first error.")
	errs.Add(newError3())
	fmt.Println(errs.FmtErrors())
	fstest.PrintTestEnd()
}
