/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: rate limiter
@author: fanky
@version: 1.0
@date: 2026-05-24
**/

package fsratelimit

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

// -------------------------------------------------------------------
// Task
// -------------------------------------------------------------------
type I_Task interface {
	Do(id int)
}

// -------------------------------------------------------------------
// Limiter
// -------------------------------------------------------------------
type S_Limiter struct {
	qsize        int
	workers      int
	chTask       chan I_Task
	once         sync.Once
	closed       atomic.Bool
	waitForDrain bool
	ctx          context.Context
	cancel       func()
}

// 创建并发限流器
// 参数：
// workers：并发数量
// qsize  ：并发队列长度，如果小于 1，会被置为 1000
// waitForDrain: 一旦关闭 Limitter，是否做完剩余任务才结束，为 true 时会等待剩余任务完成
func NewLimiter(workers int, qsize int, waitForDrain bool) *S_Limiter {
	if qsize <= 0 {
		qsize = 1000
	}
	workers = max(1, workers)
	ctx, cancel := context.WithCancel(context.Background())
	return &S_Limiter{
		qsize:        qsize,
		workers:      workers,
		chTask:       make(chan I_Task, qsize),
		waitForDrain: waitForDrain,
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (this *S_Limiter) Run(ctx context.Context) {
	go func() {
		select {
		case <-ctx.Done():
		case <-this.ctx.Done():
			// 外部 ctx 不关闭，只调用 Close 函数时，也能正常退出 goroutine
		}
		this.Close()
	}()

	work := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()
		if this.waitForDrain {
			// 关闭通道后，在队列中的任务会继续跑
			for task := range this.chTask {
				task.Do(id)
			}
			return
		}

		// 关闭通道后，所有在队列中的任务都被丢弃
		for {
			select {
			case <-ctx.Done():
				return
			case task, ok := <-this.chTask:
				if !ok {
					return
				}
				task.Do(id)
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(this.workers)
	for i := range this.workers {
		go work(&wg, i)
	}
	wg.Wait()
}

// 添加任务到任务队列
// 如果添加失败（即任务队列已经满），则返回 false。否则返回 true
func (this *S_Limiter) AddTask(task I_Task) (ok bool, err error) {
	if this.closed.Load() {
		return false, errors.New("limiter has been closed")
	}

	defer func() {
		if recover() != nil {
			err = errors.New("limiter has been closed")
		}
	}()

	select {
	case this.chTask <- task:
		return true, nil
	default:
		return false, nil
	}
}

func (this *S_Limiter) Closed() bool {
	return this.closed.Load()
}

func (this *S_Limiter) Close() {
	this.once.Do(func() {
		this.closed.Store(true)
		// 阻止新任务投递；waitForDrain 模式下同时触发 worker 排空退出
		close(this.chTask)
		this.cancel()
	})
}
