/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: functor
@author: fanky
@version: 1.0
@date: 2025-01-01
**/

package fstype

// ---------------------------------------------------------
// 单值迭代器
type F_Iter1[A any] func(func(A) bool)

// 双值迭代器
type F_Iter2[A1, A2 any] func(func(A1, A2) bool)
