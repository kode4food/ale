package builtin_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/compiler/special"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestFunctionEval(t *testing.T) {
	as := assert.New(t)
	as.EvalTo(`
		(defn say-hello ()
		  "Hello, World!")
		(say-hello)
	`, S("Hello, World!"))

	as.EvalTo(`
		(defn identity (value) value)
		(identity "foo")
	`, S("foo"))
}

func TestBadFunctionEval(t *testing.T) {
	symErr := interfaceErr("data.Integer", "data.LocalSymbol", "LocalSymbol")
	numErr := fmt.Errorf(special.UnexpectedLambdaSyntax, "99")
	vecErr := typeErr("data.Integer", "data.Vector")
	listErr := interfaceErr("data.Integer", "data.LocalSymbol", "LocalSymbol")

	as := assert.New(t)
	as.PanicWith(`(defn blah (name 99 bad) (name))`, symErr)
	as.PanicWith(`(defn blah 99 (name))`, numErr)
	as.PanicWith(`(defn 99 (x y) (+ x y))`, symErr)
	as.PanicWith(`(defn blah (99 "hello"))`, listErr)
	as.PanicWith(`(defn blah [(x) "hello"] 99)`, vecErr)
}
