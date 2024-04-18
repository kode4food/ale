package eval

import (
	"github.com/kode4food/ale/pkg/compiler/encoder"
	"github.com/kode4food/ale/pkg/compiler/generate"
	"github.com/kode4food/ale/pkg/compiler/procedure"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/read"
	"github.com/kode4food/ale/pkg/runtime"
	"github.com/kode4food/ale/pkg/runtime/isa"
)

// String evaluates the specified raw source
func String(ns env.Namespace, src data.String) data.Value {
	r := read.FromString(src)
	return Block(ns, r)
}

// Block evaluates a Sequence that a call to FromScanner might produce
func Block(ns env.Namespace, s data.Sequence) data.Value {
	var res data.Value
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		res = Value(ns, f)
	}
	return res
}

// Value evaluates the provided Value
func Value(ns env.Namespace, v data.Value) data.Value {
	defer runtime.NormalizeGoRuntimeErrors()
	e := encoder.NewEncoder(ns)
	generate.Value(e, v)
	e.Emit(isa.Return)
	return encodeAndRun(e)
}

func encodeAndRun(e encoder.Encoder) data.Value {
	encoded := e.Encode()
	fn := procedure.FromEncoded(encoded)
	closure := fn.Call().(data.Procedure)
	return closure.Call()
}
