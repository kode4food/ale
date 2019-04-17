package vm_test

import (
	"testing"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
	"gitlab.com/kode4food/ale/runtime/isa"
	"gitlab.com/kode4food/ale/runtime/vm"
)

var constants = api.Values{
	I(5),
	I(6),
	S("a thrown error"),
	api.Call(numLoopSum),
}

func makeCode(coders []isa.Coder) api.Call {
	code := make([]isa.Word, len(coders))
	for i, c := range coders {
		code[i] = c.Word()
	}
	exec := vm.NewClosure(&vm.Config{
		Code:      code,
		Constants: constants,
		StackSize: 16,
	})
	return exec(S("closure")).(api.Call)
}

func runCode(coders []isa.Coder) api.Value {
	fn := makeCode(coders)
	return fn(S("arg"))
}

func testResult(t *testing.T, res api.Value, code []isa.Coder) {
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
