package handout

import (
	"fmt"
	"sync"
)

// Control 线程开启，关闭，锁，线程数量,任务数量
type Control struct {
	wg         sync.WaitGroup
	goroutines int
	tasks      int
	ch         chan func()
}

func (C *Control) Begin(g int, tasks int) {
	if g < 10 {
		g = 10
	}
	C.goroutines = g
	C.tasks = tasks
	C.ch = make(chan func(), C.tasks)
}
func (C *Control) Gorun() {
	for range C.goroutines {
		C.wg.Add(1)
		go func() {
			defer C.wg.Done()
			for i := range C.ch {
				i()
			}
		}()
	}
}
func (C *Control) Input(f func()) {
	C.ch <- f
}
func (C *Control) Close() {
	close(C.ch)
	C.wg.Wait()
	fmt.Println("已结束")
}
