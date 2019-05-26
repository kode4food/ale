package builtin_test

import (
	"testing"

	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestFunctionEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(defn say-hello
		  []
		  "Hello, World!")
		(say-hello)
	`, S("Hello, World!"))

	as.EvalTo(`
		(defn identity [value] value)
		(identity "foo")
	`, S("foo"))
}

func TestBadFunctionEval(t *testing.T) {
	symErr := interfaceErr("data.Integer", "data.LocalSymbol", "LocalSymbol")
	listErr := typeErr("data.Integer", "*data.List")
	vecErr := typeErr("data.Integer", "data.Vector")

	as := assert.New(t)
	as.PanicWith(`(defn blah [name 99 bad] (name))`, symErr)
	as.PanicWith(`(defn blah 99 (name))`, listErr)
	as.PanicWith(`(defn 99 [x y] (+ x y))`, symErr)
	as.PanicWith(`(defn blah (99 "hello"))`, vecErr)
	as.PanicWith(`(defn blah ([x] "hello") 99)`, listErr)
}
