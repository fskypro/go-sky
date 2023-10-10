package fsworkpool

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"fsky.pro/fstest"
)

// -------------------------------------------------------------------
// 等待任务池结束
// -------------------------------------------------------------------
type s_Job struct {
	id int
}

func (this *s_Job) Do(wid int) {
	fmt.Printf("worker %d do job %d start\n", wid, this.id)
	time.Sleep(time.Second)
	fmt.Printf("worker %d do job %d end, %d jobs left\n", wid, this.id, pool.WaitingJobs())
}

// 两个工人，容纳 3 个任务的队列
var pool = NewWorkerPool(3, 22)

func TestWorkerPool(t *testing.T) {
	fstest.PrintTestBegin("test worker pool")
	defer fstest.PrintTestEnd()
	fmt.Printf("free workers: %d\n", pool.FreeWorks())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		pool.Run()
		wg.Done()
	}()

	// 添加任务
	go func() {
		for i := 1; i <= 20; i++ {
			fmt.Println("add job ", i)
			pool.Add(&s_Job{i})
		}
	}()

	// 关闭任务池
	go func() {
		time.Sleep(time.Second * 2)
		pool.Close()
		time.Sleep(time.Second)
		err := pool.Add(&s_Job{100})
		fmt.Println(11111, err)
	}()

	wg.Wait()
	fmt.Println("end")
	time.Sleep(time.Second * 2)

}

// -------------------------------------------------------------------
// 任务池不结束。只等待任务结束
// -------------------------------------------------------------------
type s_Job2 struct {
	wg *sync.WaitGroup
	id int
}

func (this *s_Job2) Do(wid int) {
	fmt.Printf("worker %d do job %d start\n", wid, this.id)
	time.Sleep(time.Second)
	fmt.Printf("worker %d do job %d end, %d jobs left\n", wid, this.id, pool2.WaitingJobs())
	this.wg.Done()
}

var pool2 = NewWorkerPool(3, 5)

func TestWorkerPool2(t *testing.T) {
	fstest.PrintTestBegin("test worker pool2")
	defer fstest.PrintTestEnd()
	fmt.Printf("free workers: %d\n", pool2.FreeWorks())

	defer pool2.Close()

	var wg sync.WaitGroup
	wg.Add(20)
	go pool2.Run()

	// 添加任务
	go func() {
		for i := 1; i <= 20; i++ {
			fmt.Println("add job ", i)
			pool2.Add(&s_Job2{&wg, i})
		}
	}()

	wg.Wait()
	fmt.Println("end")
	time.Sleep(time.Second * 2)

}
