/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: G
@author: fanky
@version: 1.0
@date: 2023-02-09
**/

package fsasyncwork

import (
	"context"
	"errors"
)

// -------------------------------------------------------------------
// job interface
// -------------------------------------------------------------------
type I_Job interface {
	Do(int)
}

// -------------------------------------------------------------------
// worker
// -------------------------------------------------------------------
type s_Worker struct {
	id int
}

// -------------------------------------------------------------------
// worker pool
// -------------------------------------------------------------------
type S_WorkerPool struct {
	chJob       chan I_Job    // 传入任务通道
	chWorker    chan s_Worker // 工作协程通道
	jobSize     int           // 任务队列长度
	workerCount int           // 工人数量

	ctx    context.Context
	cancel func()
}

// 新建任务池
// waits 为等待队列长度
func NewWorkerPool(jobs int, workers int) *S_WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	pool := &S_WorkerPool{
		chJob:       make(chan I_Job, jobs),
		chWorker:    make(chan s_Worker, workers),
		jobSize:     jobs,
		workerCount: workers,
		ctx:         ctx,
		cancel:      cancel,
	}
	for i := 1; i <= workers; i++ {
		pool.chWorker <- s_Worker{i}
	}
	return pool
}

func (this *S_WorkerPool) work(ctx context.Context, id int) {
	select {
	case <-this.ctx.Done():
		return
	case <-ctx.Done():
		return
	case job, ok := <-this.chJob:
		if !ok { return }
		job.Do(id)
	}

	func() {
		// 这里的做法是防止极低极低概率下，下面的 this.do() 退出后，使得
		// 通道已经关闭，但是这里还继续往通道里发送数据
		defer func() { recover() }()
		this.chWorker <- s_Worker{id}
	}()
}

func (this *S_WorkerPool) do(ctx context.Context) {
	if ctx == nil {
		c, cancel := context.WithCancel(context.Background())
		ctx = c
		defer cancel()
	}
	defer close(this.chJob)
	defer close(this.chWorker)
	for {
		select {
		case <-this.ctx.Done():
			return
		case <-ctx.Done():
			return
		case worker, ok := <-this.chWorker:
			// 如果有空闲工人，则取一个工人
			if !ok { return }
			// 让工人进入工作
			go this.work(ctx, worker.id)
		}
	}
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
// 添加任务
func (this *S_WorkerPool) Add(job I_Job) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("async workpoll has closed")
		}
	}()
	this.chJob <- job
	return
}

// 任务队列长度
func (this *S_WorkerPool) JobQueueSize() int {
	return this.jobSize
}

// 任务队列长度
func (this *S_WorkerPool) WaitingJobs() int {
	return len(this.chJob)
}

// 任务队列是否已满
func (this *S_WorkerPool) JobQueueIsFull() bool {
	return len(this.chJob) >= this.jobSize
}

// ---------------------------------------------------------
// 开启的协程数
func (this *S_WorkerPool) Workers() int {
	return this.workerCount
}

// 在工作的协程数
func (this *S_WorkerPool) BusyWorks() int {
	return this.workerCount - len(this.chWorker)
}

// 空闲协程数
func (this *S_WorkerPool) FreeWorks() int {
	return len(this.chWorker)
}

// ---------------------------------------------------------
// 开始进入工作
func (this *S_WorkerPool) Do() {
	this.do(nil)
}

func (this *S_WorkerPool) DoContex(ctx context.Context) {
	this.do(ctx)
}

// 结束工作
// 注意：结束工作后，不能再调用 Do/DoContex
func (this *S_WorkerPool) Close() {
	this.cancel()
}
