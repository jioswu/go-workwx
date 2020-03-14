//+build sdkcodegen

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"strings"
)

type emitter interface {
	EmitCode(spec *hir) error
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

func (e *goEmitter) EmitCode(spec *hir) error {
	e.e("// Code generated by sdkcodegen; DO NOT EDIT.\n")
	e.e("\n")
	e.e("package workwx\n")
	e.e("\n")

	for i := range spec.topics {
		err := e.emitTopic(&spec.topics[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *goEmitter) Finalize() error {
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

func (e *goEmitter) emitTopic(x *topic) error {
	for i := range x.models {
		err := e.emitModel(&x.models[i])
		if err != nil {
			return err
		}
		e.e("\n")
	}

	for i := range x.calls {
		err := e.emitAPICall(&x.calls[i])
		if err != nil {
			return err
		}
		e.e("\n")
	}

	return nil
}

func (e *goEmitter) emitModel(x *apiModel) error {
	// TODO: normalize ident according to visibility
	ident := x.ident

	e.emitDoc(ident, x.doc)
	e.e("type %s struct {\n", ident)

	for i := range x.fields {
		err := e.emitModelField(&x.fields[i])
		if err != nil {
			return err
		}
	}

	e.e("}\n")

	golangSnippets := x.inlineCodeSections["go"]
	if len(golangSnippets) > 0 {
		e.e("\n")

		for _, s := range golangSnippets {
			e.e("%s\n", s)
		}

		e.e("\n")
	}

	return nil
}

func (e *goEmitter) emitModelField(x *apiModelField) error {
	// TODO: normalize ident according to visibility
	ident := x.ident

	e.emitDoc(ident, x.doc)
	e.e("%s %s", ident, x.typ)

	if len(x.tags) > 0 {
		e.e("`")

		isFirst := true
		for k, v := range x.tags {
			if isFirst {
				isFirst = false
			} else {
				e.e(" ")
			}

			e.e("%s:\"%s\"", k, v)
		}

		e.e("`")
	}
	e.e("\n")

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

func (e *goEmitter) emitAPICall(x *apiCall) error {
	ident := x.ident

	var execMethodName string
	switch x.method {
	case apiMethodGET:
		execMethodName = "executeQyapiGet"
	case apiMethodPOSTJSON:
		execMethodName = "executeQyapiJSONPost"
	case apiMethodPOSTMedia:
		execMethodName = "executeQyapiMediaUpload"
	default:
		panic("unimplemented")
	}

	// TODO: override the receiver of method
	e.emitDoc(ident, x.doc)
	e.e("func (c *WorkwxApp) %s(req %s) (%s, error) {\n", ident, x.reqType, x.respType)
	e.e("var resp %s\n", x.respType)
	e.e("err := c.%s(\"%s\", req, &resp, %v)\n", execMethodName, x.httpURI, x.needsAccessToken)
	e.e("if err != nil {\n")
	// TODO: error_chain
	e.e("return %s{}, err\n", x.respType)
	e.e("}\n")
	e.e("if bizErr := resp.TryIntoErr(); bizErr != nil {\n")
	e.e("return %s{}, bizErr\n", x.respType)
	e.e("}\n")
	e.e("\n")
	e.e("return resp, nil\n")
	e.e("}\n")
	e.e("\n")

	return nil
}
