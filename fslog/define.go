/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: definations
@author: fanky
@version: 1.0
@date: 2024-01-11
**/

package fslog

import "strings"

type T_Level string

func Lv(lv string) T_Level {
	return T_Level(strings.ToUpper(lv))
}
