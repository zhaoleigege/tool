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

func New() *PanicG {
	return &PanicG{
		doneChan:  make(chan struct{}),
		panicChan: make(chan panicInfo),
		count:     0,
	}
}

// PanicG 处理多个协程一起运行的情况
// 这里捕获了每一个的goroutine可能出现的panic，防止整个进程结束
type PanicG struct {
	doneChan  chan struct{}
	panicChan chan panicInfo
	count     int32
}

func (pg *PanicG) Go(fc func()) *PanicG {
	atomic.AddInt32(&pg.count, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				pg.panicChan <- panicInfo{stack: debug.Stack(), source: r}
				return
			}

			pg.doneChan <- struct{}{}
		}()

		fc()
	}()

	return pg
}

func (pg *PanicG) Done() {
	for {
		select {
		case pi := <-pg.panicChan:
			fmt.Println(string(pi.stack))
		case <-pg.doneChan:
		}

		remain := atomic.AddInt32(&pg.count, -1)
		if remain == 0 {
			break
		}
	}
}
