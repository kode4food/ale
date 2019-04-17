package eval

import (
	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/compiler/generate"
	"gitlab.com/kode4food/ale/namespace"
	"gitlab.com/kode4food/ale/read"
	"gitlab.com/kode4food/ale/runtime/isa"
	"gitlab.com/kode4food/ale/runtime/vm"
)

// String evaluates the specified raw source
func String(ns namespace.Type, src api.String) api.Value {
	r := read.FromString(src)
	return Block(ns, r)
}

// Block evaluates a Sequence that a call to FromScanner might produce
func Block(ns namespace.Type, s api.Sequence) api.Value {
	var res api.Value
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		e := encoder.NewEncoder(ns)
		generate.Value(e, f)
		e.Append(isa.Return)
		res = encodeAndRun(e)
	}
	return res
}

// Value evaluates the provided Value
func Value(ns namespace.Type, v api.Value) api.Value {
	e := encoder.NewEncoder(ns)
	generate.Value(e, v)
	e.Append(isa.Return)
	return encodeAndRun(e)
}

func encodeAndRun(e encoder.Type) api.Value {
	cfg := &vm.Config{
		Globals:    e.Globals(),
		Constants:  e.Constants(),
		Code:       e.Code(),
		StackSize:  e.StackSize(),
		LocalCount: e.LocalCount(),
	}
	call := vm.NewClosure(cfg)().(api.Call)
	return call()
}
