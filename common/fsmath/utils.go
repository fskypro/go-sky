/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: utils
@author: fanky
@version: 1.0
@date: 2026-07-07
**/

package fsmath

func Abs[T T_Signed](value T) T {
	if value > 0 {
		return value
	}
	return -value
}
