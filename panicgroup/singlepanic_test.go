package panicgroup

import (
	"fmt"
	"sync"
	"testing"
)

func TestHandleGoPanic1(t *testing.T) {
	wg := &sync.WaitGroup{}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		i := i
		HandleGoPanic(func() {
			fmt.Printf("执行%d\n", i)
			wg.Done()
		})
	}

	wg.Wait()
}

func TestHandleGoPanic2(t *testing.T) {
	wg := &sync.WaitGroup{}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		i := i
		HandleGoPanic(func() {
			defer wg.Done()

			fmt.Printf("执行%d\n", i)
		})
	}

	wg.Add(1)
	HandleGoPanic(func() {
		defer wg.Done()

		panic("故意错误")
	})

	wg.Wait()
}