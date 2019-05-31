package special

import (
	"fmt"

	"gitlab.com/kode4food/ale/compiler/arity"
	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/compiler/generate"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/eval"
	"gitlab.com/kode4food/ale/namespace"
	"gitlab.com/kode4food/ale/runtime/isa"
)

// Error messages
const (
	UnpairedBindings = "bindings must be paired"
	NameAlreadyBound = "name is already bound in local scope: %s"
)

type (
	letBinding struct {
		name  data.Name
		value data.Value
	}

	letBindings []*letBinding

	uniqueNames map[data.Name]bool
)

// Eval encodes an evaluation
func Eval(e encoder.Type, args ...data.Value) {
	arity.AssertFixed(1, len(args))
	generate.Value(e, args[0])
	generate.Literal(e, evalFor(e.Globals()))
	e.Emit(isa.Call1)
}

func evalFor(ns namespace.Type) data.Call {
	return data.Call(func(args ...data.Value) data.Value {
		return eval.Value(ns, args[0])
	})
}

// Do encodes a set of expressions, returning only the final evaluation
func Do(e encoder.Type, args ...data.Value) {
	v := data.NewVector(args...)
	generate.Block(e, v)
}

// If encodes an (if cond then else) form
func If(e encoder.Type, args ...data.Value) {
	al := arity.AssertRanged(2, 3, len(args))
	generate.Branch(e,
		func() {
			generate.Value(e, args[0])
			e.Emit(isa.MakeTruthy)
		},
		func() {
			generate.Value(e, args[1])
		},
		func() {
			if al == 3 {
				generate.Value(e, args[2])
			} else {
				generate.Nil(e)
			}
		},
	)
}

// Let encodes a binding form. Binding values are evaluated first,
// and are then bound to fresh names, meaning that mutual recursion
// is not supported. For that, see `LetMutual`
func Let(e encoder.Type, args ...data.Value) {
	bindings, body := parseLet(args...)

	e.PushLocals()
	// Push the evaluated expressions to be bound
	for _, b := range bindings {
		generate.Value(e, b.value)
	}

	// Bind the popped expression results to names
	for i := len(bindings) - 1; i >= 0; i-- {
		b := bindings[i]
		e.Emit(isa.Store, e.AddLocal(b.name))
	}

	generate.Block(e, body)
	e.PopLocals()
}

func parseLet(args ...data.Value) (letBindings, data.Vector) {
	arity.AssertMinimum(2, len(args))
	b := args[0].(data.Vector)
	lb := len(b)
	if lb%2 != 0 {
		panic(fmt.Errorf(UnpairedBindings))
	}
	names := uniqueNames{}
	bindings := letBindings{}
	for i := 0; i < lb; i += 2 {
		name := b[i].(data.LocalSymbol).Name()
		names.see(name)
		value := b[i+1]
		bindings = append(bindings, newLetBinding(name, value))
	}
	return bindings, args[1:]
}

func newLetBinding(name data.Name, value data.Value) *letBinding {
	return &letBinding{
		name:  name,
		value: value,
	}
}

func (u uniqueNames) see(n data.Name) {
	if _, ok := u[n]; ok {
		panic(fmt.Errorf(NameAlreadyBound, n))
	}
	u[n] = true
}
