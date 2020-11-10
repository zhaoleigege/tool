package errors

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

type bErr struct {
	message string    // 这一次错误的信息描述
	preErr  error     // 指向上一个错误
	stack   []uintptr // 调用栈
}

// Error 实现了error接口
func (e *bErr) Error() string {
	buf := bytes.Buffer{}

	pcs := make([]uintptr, 0)
	msgs := make([]string, 0)
	for innerErr := e; innerErr != nil; {
		pcs = innerErr.stack
		msgs = append(msgs, innerErr.message)
		preErr := innerErr.preErr

		cause, ok := preErr.(*bErr)
		if !ok {
			break
		}

		innerErr = cause
	}

	for i, pc := range pcs {
		f := runtime.FuncForPC(pc - 1)
		file, line := f.FileLine(pc - 1)

		// 最外层的错误不再打印出来
		if strings.Index(file, runtime.GOROOT()) != -1 {
			break
		}

		msg := ""
		index := len(msgs) - i - 1
		if index >= 0 {
			msg = msgs[index]
		}

		buf.WriteString(fmt.Sprintf("\n\t %s:%d %s %s", file, line, f.Name(), msg))
	}

	return buf.String()
}

// Format 实现了Formatter接口，可以自定义格式化输入输出的内容
func (e *bErr) Format(s fmt.State, c rune) {
	_, _ = s.Write([]byte(e.Error()))
}

// New 返回一个包含了原始error的新error对象
// 让一开始error记录调用的堆栈信息
func New(err error, msg string) error {
	if err, ok := err.(*bErr); ok {
		return &bErr{
			message: msg,
			preErr:  err,
		}
	}

	return &bErr{
		message: msg,
		preErr:  err,
		stack:   errCallers(),
	}
}

// NewF 使用format的方式返回一个包含了原始error的新error对象
func NewF(err error, format string, a ...interface{}) error {
	msg := fmt.Sprintf(format, a...)

	if err, ok := err.(*bErr); ok {
		return &bErr{
			message: msg,
			preErr:  err,
		}
	}

	return &bErr{
		message: msg,
		preErr:  err,
		stack:   errCallers(),
	}
}

// Cause 返回原始的错误类型
func Cause(err error) error {
	for err != nil {
		cause, ok := err.(*bErr)
		if !ok {
			break
		}

		err = cause.preErr
	}

	return err
}

// errCallers 返回bErr生成的函数调用栈
func errCallers() []uintptr {
	pc := make([]uintptr, 32)
	n := runtime.Callers(3, pc)
	return pc[:n]
}
