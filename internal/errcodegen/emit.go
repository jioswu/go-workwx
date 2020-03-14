//+build sdkcodegen

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"strconv"
	"strings"
	"time"
)

type emitter interface {
	Init(retrieveTime time.Time) error
	EmitErrCode(code int64, desc string, solution string) error
	Finalize() error
}

type goEmitter struct {
	Sink io.Writer

	buf bytes.Buffer
}

var _ emitter = (*goEmitter)(nil)

func (e *goEmitter) e(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(&e.buf, format, a...)
}

func (e *goEmitter) Init(retrieveTime time.Time) error {
	e.e("// Code generated by errcodegen; DO NOT EDIT.\n")
	e.e("\n")
	e.e("package workwx\n")
	e.e("\n")
	e.e("// ErrCode 错误码类型\n")
	e.e("//\n")
	e.e("// 全局错误码文档: %s\n", errcodeDocURL)
	e.e("type ErrCode = int64\n")
	e.e("\n")
	e.e("const (\n")
	e.e("// 文档爬取时间: %s\n", retrieveTime.Format("2006-01-02 15:04:05 -0700"))
	e.e("//\n")
	e.e("// NOTE: 关于错误码的名字为何如此无聊:\n")
	e.e("//\n")
	e.e("// 官方没有给出每个错误码对应的标识符，数量太多了\n")
	e.e("// 我也懒得帮他们想，反正有文档，就先这样吧\n")
	e.e("\n")
	return nil
}

func (e *goEmitter) Finalize() error {
	e.e(")\n")

	result, err := format.Source(e.buf.Bytes())
	if err != nil {
		return err
	}

	_, err = e.Sink.Write(result)
	if err != nil {
		return err
	}

	return nil
}

func (e *goEmitter) EmitErrCode(code int64, desc string, solution string) error {
	name, ok := errcodeNameMap[code]
	if !ok {
		name = strconv.FormatInt(code, 10)
	}
	ident := fmt.Sprintf("ErrCode%s", name)
	doc := fmt.Sprintf("%s\n排查方法: %s", desc, solution)

	e.emitDoc(ident, doc)
	e.e("%s ErrCode = %d\n", ident, code)

	return nil
}

func (e *goEmitter) emitDoc(ident string, doc string) error {
	if len(doc) == 0 {
		return nil
	}

	lines := strings.Split(doc, "\n")
	e.e("// %s %s\n", ident, lines[0])

	if len(lines) == 1 {
		return nil
	}

	e.e("//\n")
	for _, l := range lines[1:] {
		e.e("// %s\n", l)
	}

	return nil
}
