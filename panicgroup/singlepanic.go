package panicgroup

import (
	"fmt"
	"runtime/debug"
)

// HandleGoPanic 对于单个的goroutine，这里可以捕获panic不至于让一个协程的panic导致整个进程的结束
// 注意闭包的传值问题需要自己解决
func HandleGoPanic(f func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("%v stack: %s", r, debug.Stack())
				return
			}
		}()

		f()
	}()
}
