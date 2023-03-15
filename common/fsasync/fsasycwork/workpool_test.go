package fsasyncwork

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"fsky.pro/fstest"
)

type s_Job struct {
	id int
}

func (this *s_Job) Do(wid int) {
	fmt.Printf("worker %d do job %d start\n", wid, this.id)
	time.Sleep(time.Second)
	fmt.Printf("worker %d do job %d end, %d jobs left\n", wid, this.id, pool.WaitingJobs())
}

var pool = NewWorkerPool(100, 2)

func TestWorkerPool(t *testing.T) {
	fstest.PrintTestBegin("test worker pool")
	defer fstest.PrintTestEnd()
	fmt.Printf("free workers: %d\n", pool.FreeWorks())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		pool.Do()
		wg.Done()
	}()

	for i := 1; i <= 20; i++ {
		pool.Add(&s_Job{i})
	}

	go func() {
		time.Sleep(time.Second * 5)
		pool.Close()
		time.Sleep(time.Second)
		err := pool.Add(&s_Job{100})
		fmt.Println(11111, err)
	}()

	wg.Wait()
	time.Sleep(time.Second * 2)

}
