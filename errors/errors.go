package errors

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

type bErr struct {
	message string    // 这一次错误的信息描述
	preErr  *bErr     // 指向上一个错误
	stack   []uintptr // 调用栈
}

// Error 实现了error接口
func (e *bErr) Error() string {

	buf := bytes.Buffer{}

	pcs := make([]uintptr, 0)
	msgs := make([]string, 0)
	for innerErr := e; innerErr != nil; innerErr = innerErr.preErr {
		pcs = innerErr.stack
		msgs = append(msgs, innerErr.message)
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

		buf.WriteString(fmt.Sprintf("\t %s:%d %s %s\n", file, line, f.Name(), msg))
	}

	return buf.String()
}

// New 返回一个包含了原始error的新error对象
func New(err error, msg string) error {
	if err, ok := err.(*bErr); ok {
		return &bErr{
			message: msg,
			preErr:  err,
		}
	}

	return &bErr{
		message: msg,
		preErr:  nil,
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
		preErr:  nil,
		stack:   errCallers(),
	}
}

// errCallers 返回bErr生成的函数调用栈
func errCallers() []uintptr {
	pc := make([]uintptr, 32)
	n := runtime.Callers(3, pc)
	return pc[:n]
}
