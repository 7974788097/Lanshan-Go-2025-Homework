package handout

import (
	"fmt"
	"sync"
)

type Control struct {
	Lock       sync.Mutex
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
			C.Lock.Lock()
			defer C.wg.Done()
			for i := range C.ch {
				i()
			}
			C.Lock.Unlock()
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
