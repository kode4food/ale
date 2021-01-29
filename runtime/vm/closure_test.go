package vm_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/runtime/isa"
	"github.com/kode4food/ale/runtime/vm"
)

var constants = data.Values{
	I(5),
	I(6),
	S("a thrown error"),
	data.Applicative(numLoopSum),
}

func makeCode(coders []isa.Coder) data.Function {
	code := make([]isa.Word, len(coders))
	for i, c := range coders {
		code[i] = c.Word()
	}
	lambda := &vm.Lambda{
		Code:      code,
		Constants: constants,
		StackSize: 16,
	}
	closure := lambda.Call(S("closure"))
	return closure.(data.Function)
}

func runCode(coders []isa.Coder) data.Value {
	code := makeCode(coders)
	return code.Call(S("arg"))
}

func testResult(t *testing.T, res data.Value, code []isa.Coder) {
	as := assert.New(t)
	r := runCode(code)
	as.Equal(res, r)
}

func testPanic(t *testing.T, errStr string, code []isa.Coder) {
	as := assert.New(t)
	defer as.ExpectPanic(errStr)
	runCode(code)
}

func TestSimple(t *testing.T) {
	testResult(t, I(11), []isa.Coder{
		isa.Const, isa.Index(0),
		isa.Const, isa.Index(1),
		isa.Add,
		isa.Return,
	})

	testResult(t, I(0), []isa.Coder{
		isa.One,
		isa.NegOne,
		isa.Add,
		isa.Return,
	})

	testResult(t, I(0), []isa.Coder{
		isa.Zero,
		isa.Const, isa.Index(0),
		isa.Mul,
		isa.Return,
	})

	testResult(t, I(-1), []isa.Coder{
		isa.Zero,
		isa.One,
		isa.Sub,
		isa.Return,
	})

	testResult(t, S("closure"), []isa.Coder{
		isa.Closure, isa.Index(0),
		isa.Return,
	})

	testResult(t, S("arg"), []isa.Coder{
		isa.Arg, isa.Index(0),
		isa.Return,
	})
}

func TestCalls(t *testing.T) {
	testResult(t, I(17), []isa.Coder{
		isa.Const, isa.Index(0),
		isa.Const, isa.Index(0),
		isa.Const, isa.Index(1),
		isa.One,
		isa.Const, isa.Index(3),
		isa.Call, isa.Count(3),
		isa.Add,
		isa.Return,
	})
}

func TestErrors(t *testing.T) {
	testPanic(t, "a thrown error", []isa.Coder{
		isa.Const,
		isa.Index(2),
		isa.Panic,
	})
}

func TestExplosions(t *testing.T) {
	testPanic(t, "runtime error: index out of range", []isa.Coder{
		isa.Return,
	})
}
