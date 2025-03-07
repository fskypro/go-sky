/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 安全协程
@author: fanky
@version: 1.0
@date: 2024-09-18
**/

// ------------------------------------------------------------------
// 协程崩溃后，不会造成主程序退出
// 使用方法：
//   go SafeGo1[int](
//     func(err any) {
//       log.Println("协程崩溃")
//     },
//     func(i int) {
//       panic("作为例子，这里触发崩溃")
//     },
//   )
// 说明:
//   SafeGo 的第一个参数 ferr 为崩溃调用函数，如果调用第二个参数产生了崩溃，则 ferr 会被调用，
//   ferr 的参数为崩溃异常对象，它通常是一个 error 对象
// ------------------------------------------------------------------

package fsky

// 不带参数协程
func SafeGo0(ferr func(any), f func()) {
	defer func() {
		if err := recover(); err != nil {
			ferr(err)
		}
	}()
	f()
}

// 不带一个参数的协程
func SafeGo1[T1 any](ferr func(any), f func(T1), a1 T1) {
	defer func() {
		if err := recover(); err != nil {
			ferr(err)
		}
	}()
	f(a1)
}

// 不带两个参数的协程
func SafeGo2[T1, T2 any](ferr func(any), f func(T1, T2), a1 T1, a2 T2) {
	defer func() {
		if err := recover(); err != nil {
			ferr(err)
		}
	}()
	f(a1, a2)
}

// 带三个参数的协程
func SafeGo3[T1, T2, T3 any](ferr func(any), f func(T1, T2, T3), a1 T1, a2 T2, a3 T3) {
	defer func() {
		if err := recover(); err != nil {
			ferr(err)
		}
	}()
	f(a1, a2, a3)
}

// 带四个参数的协程
func SafeGo4[T1, T2, T3, T4 any](ferr func(any), f func(T1, T2, T3, T4), a1 T1, a2 T2, a3 T3, a4 T4) {
	defer func() {
		if err := recover(); err != nil {
			ferr(err)
		}
	}()
	f(a1, a2, a3, a4)
}

// 带五个参数的协程
func SafeGo5[T1, T2, T3, T4, T5 any](ferr func(any), f func(T1, T2, T3, T4, T5), a1 T1, a2 T2, a3 T3, a4 T4, a5 T5) {
	defer func() {
		if err := recover(); err != nil {
			ferr(err)
		}
	}()
	f(a1, a2, a3, a4, a5)
}

// 带六个参数的协程
func SafeGo6[T1, T2, T3, T4, T5, T6 any](ferr func(any), f func(T1, T2, T3, T4, T5, T6), a1 T1, a2 T2, a3 T3, a4 T4, a5 T5, a6 T6) {
	defer func() {
		if err := recover(); err != nil {
			ferr(err)
		}
	}()
	f(a1, a2, a3, a4, a5, a6)
}

// 带七个参数的协程
func SafeGo7[T1, T2, T3, T4, T5, T6, T7 any](ferr func(any), f func(T1, T2, T3, T4, T5, T6, T7), a1 T1, a2 T2, a3 T3, a4 T4, a5 T5, a6 T6, a7 T7) {
	defer func() {
		if err := recover(); err != nil {
			ferr(err)
		}
	}()
	f(a1, a2, a3, a4, a5, a6, a7)
}

// 带八个参数的协程
func SafeGo8[T1, T2, T3, T4, T5, T6, T7, T8 any](ferr func(any), f func(T1, T2, T3, T4, T5, T6, T7, T8), a1 T1, a2 T2, a3 T3, a4 T4, a5 T5, a6 T6, a7 T7, a8 T8) {
	defer func() {
		if err := recover(); err != nil {
			ferr(err)
		}
	}()
	f(a1, a2, a3, a4, a5, a6, a7, a8)
}

// ------------------------------------------------------------------
// 以协程对象的方式实现安全协程
// 使用方法：
//
//	 type Worker struct {
//	    a int
//	    s string
//	 }
//	 func(self Wroker) Do() {
//	   fmt.Println(100 / this.a)
//	 }
//	 func(self Worker) Error(err any) {
//	   fmt.Println("panic:", err)
//	 }
//
//	go SafeGo(Worker{0, "str"})
//
// ------------------------------------------------------------------
type I_SafeGoWorker interface {
	Do()
	Error(any)
}

// 调用 worker 中的 Do 方法来执行异步操作，如果 Do 方法产生 panic，则 worker 的 Error 方法会被调用
func SafeGo(worker I_SafeGoWorker) {
	defer func() {
		if err := recover(); err != nil {
			worker.Error(err)
		}
	}()
	worker.Do()
}
