package panicgroup

import (
	"fmt"
	"testing"
	"time"
)

func TestPanicG_Done(t *testing.T) {
	f1 := func() {
		time.Sleep(1 * time.Second)
		panic("故意的错误")
	}

	f2 := func() {
		time.Sleep(3 * time.Second)
		fmt.Println("这里没有问题")
	}

	pg := New().Go(f1).Go(f2)
	fmt.Println("等待结束")
	pg.Done()
}
