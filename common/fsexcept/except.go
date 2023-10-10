/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 模拟异常
@author: fanky
@version: 1.0
@date: 2022-09-07
**/

package fsexcept

type S_Except struct {
	Try     func(S_Except)
	Catch   func(any)
	Finally func()
}

func (self S_Except) Do() {
	if self.Finally != nil {
		defer self.Finally()
	}
	if self.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				self.Catch(r)
			}
		}()
	}
	self.Try(self)
}

func Throw(obj any) {
	panic(obj)
}
