package main

import (
	"fmt"
	"github.com/zhaoleigege/tool/errors"
	"sync"
)

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
	switch e := errors.Cause(err).(type) {
	case *NameErr:
		fmt.Printf("NameErr: %s: %d\n", e.Name, e.Age)
	default:
		switch e {
		case errors.NotExist:
			fmt.Printf("原始错误是NotExist")
		}
	}

	fmt.Printf("最里面的错误: %v\n", errors.Cause(err))
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
		//err := errors.New(errors.NotExist, "bar发生了错误")
		err := errors.New(&NameErr{
			Name: "test",
			Age:  20,
		}, "bar发生了错误")
		return errors.New(err, "闭包函数错误")
	}()
}
