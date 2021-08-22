/**
@copyright: fantasysky 2016
@brief: 实现定时器
@author: fanky
@version: 1.0
@date: 2019-04-03
**/

package fstime

import (
	"time"

	"golang.org/x/net/context"
)

// -------------------------------------------------------------------
// Timer
// -------------------------------------------------------------------
type S_Timer struct {
	timer  *time.Timer
	ctx    context.Context
	cancel context.CancelFunc
	C      chan bool
}

func NewTimer(d time.Duration) *S_Timer {
	ctx, cancel := context.WithCancel(context.Background())
	this := &S_Timer{
		timer:  time.NewTimer(d),
		ctx:    ctx,
		cancel: cancel,
		C:      make(chan bool),
	}
	go func() {
		defer func() { this.timer.Stop() }()
		select {
		case <-this.timer.C:
			this.C <- true
		case <-this.ctx.Done():

			this.C <- false
		}
	}()
	return this
}

// 取消定时器
func (this *S_Timer) Cancel() {
	this.cancel()
}

// 重新设置延时
// 注意：如果一个 timer 已经结束，则再对其 Reset 是无效的
func (this *S_Timer) Reset(d time.Duration) bool {
	return this.timer.Reset(d)
}

// -------------------------------------------------------------------
// TimerFunc
// -------------------------------------------------------------------
// timer 回调函数
// 第一个参数为 true 表示 timer 正常结束，若为 false，则表示 timer 被取消
type F_TimerFunc func(bool, ...interface{})

type S_TimerFunc struct {
	timer   *time.Timer
	ctx     context.Context
	cancel  context.CancelFunc
	fun     F_TimerFunc
	funArgs []interface{}
}

func NewTimerFunc(d time.Duration, fun F_TimerFunc, funArgs ...interface{}) *S_TimerFunc {
	ctx, cancel := context.WithCancel(context.Background())
	this := &S_TimerFunc{
		timer:   time.NewTimer(d),
		ctx:     ctx,
		cancel:  cancel,
		fun:     fun,
		funArgs: funArgs,
	}
	go func() {
		defer func() { this.timer.Stop() }()
		select {
		case <-this.timer.C:
			this.onTimerEnd(true)
		case <-this.ctx.Done():
			this.onTimerEnd(false)
		}
	}()
	return this
}

func (this *S_TimerFunc) onTimerEnd(arrived bool) {
	if this.fun != nil {
		this.fun(arrived, this.funArgs...)
	}
}

// 取消定时器
func (this *S_TimerFunc) Cancel() {
	this.cancel()
}

// 重新设置延时
// 注意：如果一个 timer 已经结束，则再对其 Reset 是无效的
func (this *S_TimerFunc) Reset(d time.Duration) bool {
	return this.timer.Reset(d)
}
