package main

import (
	"errors"
	"fmt"
	"github.com/zhaoleigege/tool/errorx"
	"sync"
)

var NotExist = errors.New("not exist")

type NameErr struct {
	Name string
	Age  int64
}

func (e *NameErr) Error() string {
	return fmt.Sprintf("%s:%d", e.Name, e.Age)
}

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
	fmt.Printf("出现错误: %s\n", err)
	fmt.Printf("出现错误: %v\n", err)
	switch e := errorx.Cause(err).(type) {
	case *NameErr:
		fmt.Printf("NameErr: %s: %d\n", e.Name, e.Age)
	default:
		switch e {
		case NotExist:
			fmt.Printf("原始错误是NotExist")
		}
	}

	fmt.Printf("最里面的错误: %v\n", errorx.Cause(err))
}

func foo() error {
	fmt.Println("foo")
	if err := bar(); err != nil {
		return errorx.New(err, "foo出现错误")
	}

	return nil
}

func bar() error {
	fmt.Println("bar")
	return func() error {
		//err := errors.New(errors.NotExist, "bar发生了错误")
		err := errorx.New(&NameErr{
			Name: "test",
			Age:  20,
		}, "bar发生了错误")
		return errorx.New(err, "闭包函数错误")
	}()
}
