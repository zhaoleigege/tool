package main

import (
	"fmt"
	"github.com/zhaoleigege/tool/errors"
	"sync"
)

func main() {
	Err2()
}

func Err2() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		Err1()
		wg.Done()
	}()

	wg.Wait()
}

func Err1() {
	err := foo()
	fmt.Printf("出现错误: %s", err)
}

func foo() error {
	fmt.Println("foo")
	if err := bar(); err != nil {
		return errors.New(err, "foo出现错误")
	}

	return nil
}

func bar() error {
	fmt.Println("bar")
	return func() error {
		err := errors.New(nil, "bar发生了错误")
		return errors.New(err, "闭包函数错误")
	}()
}
