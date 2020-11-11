package errorx

import (
	"errors"
	"fmt"
	"testing"
)

func Example_new() {
	err := New(nil, "错误")
	fmt.Println(err)
	// Output:
	// /Users/klook/Desktop/tool/errorx/example_test.go:10 errorx.Example_new 错误
}

func Example_newf() {
	err := NewF(nil, "打印%s", "错误")
	fmt.Println(err)
	// Output:
	// /Users/klook/Desktop/tool/errorx/example_test.go:17 errorx.Example_newf 打印错误
}

func Example_multipl_error() {
	err := foo()
	fmt.Println(err)
	// Output:
	//	 /Users/klook/Desktop/tool/errorx/example_test.go:44 errorx.bar.func1 闭包函数错误, bar发生了错误
	//	 /Users/klook/Desktop/tool/errorx/example_test.go:45 errorx.bar foo出现错误
	//	 /Users/klook/Desktop/tool/errorx/example_test.go:34 errorx.foo
	//	 /Users/klook/Desktop/tool/errorx/example_test.go:24 errorx.Example_multipl_error
}

func foo() error {
	if err := bar(); err != nil {
		return New(err, "foo出现错误")
	}

	return nil
}

func bar() error {
	return func() error {
		err := errors.New("bar发生了错误")
		return New(err, "闭包函数错误")
	}()
}

type NameErr struct {
	Name string
	Age  int64
}

func (e *NameErr) Error() string {
	return fmt.Sprintf("%s:%d", e.Name, e.Age)
}

func TestCause(t *testing.T) {
	err1 := &NameErr{
		Name: "test",
		Age:  10,
	}

	err2 := New(err1, "包装err1错误")

	switch e := Cause(err2).(type) {
	case *NameErr:
	default:
		t.Errorf("包装类型不符合预期: %s", e)
	}
}

var NotExist = errors.New("not exist")

func TestCause2(t *testing.T) {
	err2 := New(NotExist, "包装err1错误")

	switch e := Cause(err2).(type) {
	case *NameErr:
	default:
		switch e {
		case NotExist:
		default:
			t.Errorf("包装类型不符合预期: %s", e)
		}
	}
}
