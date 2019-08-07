/**
@copyright: fantasysky 2016
@brief: 实现定时器
@author: fanky
@version: 1.0
@date: 2019-04-03
**/

package fstime

import "time"
import "context"

// F_TimerFunc 定时器回调函数
type F_TimerFunc func(time.Time, interface{})

// F_TimerEndFunc 定时器结束回调函数
type F_TimerEndFunc func(interface{})

// Timer 定时器
// 注意：
//	1、当调用定时器的 Stop 后，还可以通过 Start 再次启动
//	2、定时器回调是在另一个协程里，所以对于共享数据，要加锁
//	3、并不能保证最后一个间隔回调一定会比 EndFunc 更早地调用。而事实上，EndFunc 通常会比最后一个间隔回调更早调用
type Timer struct {
	ctx      context.Context
	ctxFunc  context.CancelFunc
	interval time.Duration
	cb       F_TimerFunc

	nowExecs  int            // 当前执行次数
	ExecTimes int            // 最多执行次数
	EndTime   time.Time      // 指定结束时间
	EndFunc   F_TimerEndFunc // timer 结束时回调函数
}

func NewTimer(interval time.Duration, cb F_TimerFunc) *Timer {
	return &Timer{
		interval: interval,
		cb:       cb,

		nowExecs:  0,
		ExecTimes: 0,
		EndTime:   time.Time{},
		EndFunc:   nil,
	}
}

func (this *Timer) Start(start time.Duration, data interface{}) {
	this.ctx, this.ctxFunc = context.WithCancel(context.Background())
	this.nowExecs = 0

	ch := time.After(this.interval)
	var handle = func() {
	L:
		for {
			select {
			case <-this.ctx.Done():
				break L
			case now := <-ch:
				if this.cb != nil {
					go this.cb(now, data)
				}
				if this.ExecTimes > 0 {
					this.nowExecs += 1
					if this.nowExecs >= this.ExecTimes {
						break L
					} else {
						ch = time.After(this.interval)
					}
				} else if !this.EndTime.IsZero() && now.After(this.EndTime) {
					break L
				} else {
					ch = time.After(this.interval)
				}
			}
		}
		if this.EndFunc != nil {
			go this.EndFunc(data)
		}
	}
	go handle()
	ch = time.After(start)
}

func (this *Timer) Stop() {
	// 如果还没 Star 就调用 Stop，则 this.ctxFunc 为 nil
	if this.ctxFunc != nil {
		this.ctxFunc()
		this.ctxFunc = nil
	}
}
