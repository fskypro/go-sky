/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 初始化器
@author: fanky
@version: 1.0
@date: 2022-03-19
**/

package fsinit

type E_Level int

// L1 优先于 L5 调用
const (
	L1 E_Level = 1
	L2         = 2
	L3         = 3
	L4         = 4
	L5         = 5
)

var initers map[E_Level][]func() = map[E_Level][]func(){}

func Add(lv E_Level, initer func()) {
	l := initers[lv]
	if l == nil {
		l = []func(){}
	}
	l = append(l, initer)
	initers[lv] = l
}

func Init() {
	for l := L1; l <= L5; l++ {
		inits := initers[l]
		if inits == nil {
			continue
		}
		for _, initer := range inits {
			initer()
		}
	}
}

func InitLevel(l E_Level) {
	inits := initers[l]
	if inits == nil {
		return
	}
	for _, initer := range inits {
		initer()
	}
}
