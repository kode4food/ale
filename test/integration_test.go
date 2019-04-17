package test

import (
	"fmt"
	"testing"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/bootstrap"
	"gitlab.com/kode4food/ale/eval"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
	"gitlab.com/kode4food/ale/namespace"
)

var (
	manager = namespace.NewManager()
	ready   bool
)

func interfaceErr(concrete, intf, method string) error {
	err := "interface conversion: %s is not %s: missing method %s"
	return fmt.Errorf(fmt.Sprintf(err, concrete, intf, method))
}

func typeErr(concrete, expected string) error {
	err := "interface conversion: api.Value is %s, not %s"
	return fmt.Errorf(fmt.Sprintf(err, concrete, expected))
}

func runCode(src string) api.Value {
	if !ready {
		bootstrap.Into(manager)
		ready = true
	}
	ns := manager.GetAnonymous()
	return eval.String(ns, api.String(src))
}

func testCode(t *testing.T, src string, expect api.Value) {
	as := assert.New(t)
	as.Equal(expect, runCode(src))
}

func testBadCode(t *testing.T, src string, err error) {
	as := assert.New(t)

	defer as.ExpectPanic(err.Error())
	runCode(src)
}

func TestDo(t *testing.T) {
	testCode(t, `
		(do
			55
			(if true 99 33))
	`, F(99))
}

func TestTrueFalse(t *testing.T) {
	testCode(t, `true`, api.True)
	testCode(t, `false`, api.False)
	testCode(t, `nil`, api.Nil)
}

func TestReadEval(t *testing.T) {
	testCode(t, `
		(eval (read "(str \"hello\" \"you\" \"test\")"))
	`, S("helloyoutest"))
}
