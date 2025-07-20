package optimize_test

import (
	"testing"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/compiler/encoder"
	"github.com/kode4food/ale/internal/compiler/generate"
	"github.com/kode4food/ale/internal/compiler/procedure"
	"github.com/kode4food/ale/internal/runtime/isa"
	"github.com/kode4food/ale/internal/runtime/vm"
)

func mustFromEncoded(e *encoder.Encoded) *vm.Procedure {
	res, err := procedure.FromEncoded(e)
	if err != nil {
		panic(err)
	}
	return res
}

func getInlineTestNamespace() env.Namespace {
	e := env.NewEnvironment()
	ns := e.GetRoot()

	bind := func(n data.Local, fn func(enc encoder.Encoder), a ...ale.Value) {
		enc := encoder.NewEncoder(ns)
		fn(enc)
		proc := mustFromEncoded(enc.Encode())
		_ = env.BindPublic(ns, n, proc.Call(a...))
	}

	// only uses args, no jumps
	bind("+", func(enc encoder.Encoder) {
		enc.Emit(isa.Arg, 0)
		enc.Emit(isa.Arg, 1)
		enc.Emit(isa.Add)
		enc.Emit(isa.Return)
	})

	// uses an arg, and draws from captured, no jumps
	bind("+6", func(enc encoder.Encoder) {
		enc.Emit(isa.Arg, 0)
		enc.Emit(isa.Closure, 0)
		enc.Emit(isa.Add)
		enc.Emit(isa.Return)
	}, I(6))

	// uses args, manipulates locals, no jumps
	bind("+l", func(enc encoder.Encoder) {
		enc.Emit(isa.Arg, 0)
		enc.Emit(isa.Store, 0)
		enc.Emit(isa.Arg, 1)
		enc.Emit(isa.Store, 1)
		enc.Emit(isa.Load, 0)
		enc.Emit(isa.Load, 1)
		enc.Emit(isa.Add)
		enc.Emit(isa.Return)
	})

	bind("+cl", func(enc encoder.Encoder) {
		enc.Emit(isa.Arg, 0)
		enc.Emit(isa.Const, enc.AddConstant(env.MustResolveValue(ns, LS("+6"))))
		enc.Emit(isa.Call1)
		enc.Emit(isa.Store, 0)
		enc.Emit(isa.Const, enc.AddConstant(I(8)))
		enc.Emit(isa.Load, 0)
		enc.Emit(isa.Const, enc.AddConstant(env.MustResolveValue(ns, LS("+l"))))
		enc.Emit(isa.Call, 2)
		enc.Emit(isa.Zero)
		enc.Emit(isa.Add)
		enc.Emit(isa.Return)
	}, I(6))

	bind("diff", func(enc encoder.Encoder) {
		_ = generate.Branch(enc, func(encoder.Encoder) error {
			enc.Emit(isa.Arg, 0)
			enc.Emit(isa.Dup)
			enc.Emit(isa.Store, 0)
			enc.Emit(isa.Arg, 1)
			enc.Emit(isa.Dup)
			enc.Emit(isa.Store, 1)
			enc.Emit(isa.NumLt)
			return nil
		}, func(encoder.Encoder) error {
			enc.Emit(isa.Load, 1)
			enc.Emit(isa.Load, 0)
			return nil
		}, func(encoder.Encoder) error {
			enc.Emit(isa.Load, 0)
			enc.Emit(isa.Load, 1)
			return nil
		})
		enc.Emit(isa.Sub)
		enc.Emit(isa.Return)
	})
	return ns
}

func TestInlineArgStacking(t *testing.T) {
	as := assert.New(t)

	ns := getInlineTestNamespace()
	enc := encoder.NewEncoder(ns)
	enc.Emit(isa.Const, enc.AddConstant(I(6)))
	enc.Emit(isa.Const, enc.AddConstant(I(8)))
	enc.Emit(isa.Const, enc.AddConstant(env.MustResolveValue(ns, LS("+"))))
	enc.Emit(isa.Call, 2)
	enc.Emit(isa.Const, enc.AddConstant(I(8)))
	enc.Emit(isa.Add)
	enc.Emit(isa.Return)

	call := mustFromEncoded(enc.Encode()).Call().(*vm.Closure)
	as.Equal(I(22), call.Call())
	as.Equal(2, len(call.Constants))
	as.Equal(I(6), call.Constants[0])
	as.Equal(I(8), call.Constants[1])
	as.Instructions(isa.Instructions{
		isa.Const.New(0),
		isa.Const.New(1),
		isa.Store.New(0),
		isa.Store.New(1),
		isa.Load.New(0),
		isa.Load.New(1),
		isa.Add.New(),
		isa.Const.New(1),
		isa.Add.New(),
		isa.Return.New(),
	}, call.Code)
}

func TestInlineArgClosure(t *testing.T) {
	as := assert.New(t)

	ns := getInlineTestNamespace()
	enc := encoder.NewEncoder(ns)
	enc.Emit(isa.Const, enc.AddConstant(I(8)))
	enc.Emit(isa.Const, enc.AddConstant(env.MustResolveValue(ns, LS("+6"))))
	enc.Emit(isa.Call1)
	enc.Emit(isa.Const, enc.AddConstant(I(8)))
	enc.Emit(isa.Add)
	enc.Emit(isa.Return)

	call := mustFromEncoded(enc.Encode()).Call().(*vm.Closure)
	as.Equal(I(22), call.Call())
	as.Equal(2, len(call.Constants))
	as.Equal(I(8), call.Constants[0])
	as.Equal(I(6), call.Constants[1])
	as.Instructions(isa.Instructions{
		isa.Const.New(0),
		isa.Const.New(1),
		isa.Add.New(),
		isa.Const.New(0),
		isa.Add.New(),
		isa.Return.New(),
	}, call.Code)
}

func TestInlineLocals(t *testing.T) {
	as := assert.New(t)

	ns := getInlineTestNamespace()
	enc := encoder.NewEncoder(ns)
	enc.Emit(isa.Const, enc.AddConstant(I(6)))
	enc.Emit(isa.Const, enc.AddConstant(I(8)))
	enc.Emit(isa.Const, enc.AddConstant(env.MustResolveValue(ns, LS("+l"))))
	enc.Emit(isa.Call, 2)
	enc.Emit(isa.Const, enc.AddConstant(I(8)))
	enc.Emit(isa.Add)
	enc.Emit(isa.Return)

	call := mustFromEncoded(enc.Encode()).Call().(*vm.Closure)
	as.Equal(I(22), call.Call())
	as.Equal(2, len(call.Constants))
	as.Equal(I(6), call.Constants[0])
	as.Equal(I(8), call.Constants[1])
	as.Instructions(isa.Instructions{
		isa.Const.New(0),
		isa.Const.New(1),
		isa.Store.New(0),
		isa.Store.New(1),
		isa.Load.New(0),
		isa.Load.New(1),
		isa.Add.New(),
		isa.Const.New(1),
		isa.Add.New(),
		isa.Return.New(),
	}, call.Code)
}

func TestInlineNestedLocals(t *testing.T) {
	as := assert.New(t)

	ns := getInlineTestNamespace()
	enc := encoder.NewEncoder(ns)
	enc.Emit(isa.Const, enc.AddConstant(I(6)))
	enc.Emit(isa.Const, enc.AddConstant(I(8)))
	enc.Emit(isa.Const, enc.AddConstant(env.MustResolveValue(ns, LS("+cl"))))
	enc.Emit(isa.Call, 2)
	enc.Emit(isa.Const, enc.AddConstant(I(8)))
	enc.Emit(isa.Add)
	enc.Emit(isa.Return)

	call := mustFromEncoded(enc.Encode()).Call().(*vm.Closure)
	as.Equal(I(30), call.Call())
	as.Equal(2, len(call.Constants))
	as.Equal(I(6), call.Constants[0])
	as.Equal(I(8), call.Constants[1])
	as.Instructions(isa.Instructions{
		isa.Const.New(0),
		isa.Const.New(1),
		isa.Store.New(0),
		isa.Pop.New(),
		isa.Load.New(0),
		isa.Const.New(0),
		isa.Add.New(),
		isa.Store.New(1),
		isa.Const.New(1),
		isa.Store.New(2),
		isa.Load.New(1),
		isa.Load.New(2),
		isa.Add.New(),
		isa.Zero.New(),
		isa.Add.New(),
		isa.Const.New(1),
		isa.Add.New(),
		isa.Return.New(),
	}, call.Code)
}

func TestDiff(t *testing.T) {
	as := assert.New(t)

	ns := getInlineTestNamespace()
	diff := env.MustResolveValue(ns, LS("diff")).(*vm.Closure)
	as.Equal(I(2), diff.Call(I(5), I(7)))
	as.Equal(I(2), diff.Call(I(7), I(5)))

	enc := encoder.NewEncoder(ns)
	as.NoError(generate.Branch(enc, func(encoder.Encoder) error {
		enc.Emit(isa.Arg, 0)
		enc.Emit(isa.Arg, 1)
		enc.Emit(isa.NumEq)
		return nil
	}, func(encoder.Encoder) error {
		enc.Emit(isa.Zero)
		return nil
	}, func(encoder.Encoder) error {
		enc.Emit(isa.Arg, 0)
		enc.Emit(isa.Arg, 1)
		enc.Emit(isa.Const, enc.AddConstant(diff))
		enc.Emit(isa.Call, 2)
		return nil
	}))
	enc.Emit(isa.PosInt, 1)
	enc.Emit(isa.Add)
	enc.Emit(isa.Return)

	call := mustFromEncoded(enc.Encode()).Call().(*vm.Closure)
	as.Equal(I(1), call.Call(I(5), I(5)))
	as.Equal(I(2), call.Call(I(5), I(6)))
	as.Equal(I(2), call.Call(I(6), I(5)))
	as.Equal(I(5), call.Call(I(2), I(6)))
	as.Equal(I(5), call.Call(I(6), I(2)))
	as.Instructions(isa.Instructions{
		isa.Arg.New(0),
		isa.Arg.New(1),
		isa.NumEq.New(),
		isa.CondJump.New(23),
		isa.Arg.New(0),
		isa.Arg.New(1),
		isa.Store.New(0),
		isa.Store.New(1),
		isa.Load.New(0),
		isa.Dup.New(),
		isa.Store.New(2),
		isa.Load.New(1),
		isa.Dup.New(),
		isa.Store.New(3),
		isa.NumLt.New(),
		isa.CondJump.New(19),
		isa.Load.New(2),
		isa.Load.New(3),
		isa.Jump.New(21),
		isa.Load.New(3),
		isa.Load.New(2),
		isa.Sub.New(),
		isa.Jump.New(24),
		isa.Zero.New(),
		isa.PosInt.New(1),
		isa.Add.New(),
		isa.Return.New(),
	}, call.Code)
}
