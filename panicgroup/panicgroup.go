package panicgroup

import (
	"fmt"
	"runtime/debug"
	"sync/atomic"
)

type panicInfo struct {
	stack  []byte
	source interface{}
}

func (pi panicInfo) Format(s fmt.State, c rune) {
	_, _ = s.Write([]byte(pi.Error()))
}

func (pi panicInfo) Error() string {
	return fmt.Sprintf("r: %v, stack: %s", pi.source, string(pi.stack))
}

func New() *PanicG {
	return &PanicG{
		stopChan:  make(chan struct{}),
		panicChan: make(chan panicInfo),
		doneChan:  make(chan struct{}),
		count:     0,
	}
}

// PanicG 处理多个协程一起运行的情况
// 这里捕获了每一个的goroutine可能出现的panic，防止整个进程结束
type PanicG struct {
	stopChan  chan struct{}
	panicChan chan panicInfo
	doneChan  chan struct{}
	count     int32
}

func (pg *PanicG) Go(fc func()) *PanicG {
	atomic.AddInt32(&pg.count, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				select {
				case pg.panicChan <- panicInfo{stack: debug.Stack(), source: r}:
				case <-pg.stopChan:
				}
				return
			}

			select {
			case pg.doneChan <- struct{}{}:
			case <-pg.stopChan:
			}
		}()

		fc()
	}()

	return pg
}

func (pg *PanicG) Done() error {
	for {
		select {
		case pi := <-pg.panicChan:
			close(pg.stopChan)
			return pi
		case <-pg.doneChan:
			if atomic.AddInt32(&pg.count, -1) == 0 {
				close(pg.stopChan)
				return nil
			}
		}
	}

	return nil
}
