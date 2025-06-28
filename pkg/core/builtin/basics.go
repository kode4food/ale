package builtin

import (
	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/internal/runtime"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/read"
)

var emptyNamespace = env.NewEnvironment().GetRoot()

// Recover invokes a function and runs a recovery function if Go panics
var Recover = data.MakeProcedure(func(args ...data.Value) (res data.Value) {
	body := args[0].(data.Procedure)
	rescue := args[1].(data.Procedure)

	defer func() {
		if rec := recover(); rec != nil {
			switch rec := runtime.NormalizeGoRuntimeError(rec).(type) {
			case data.Value:
				res = rescue.Call(rec)
			case error:
				res = rescue.Call(data.String(rec.Error()))
			default:
				panic(debug.ProgrammerError("recover returned invalid result"))
			}
		}
	}()

	return body.Call()
}, 2)

// Defer invokes a cleanup function, no matter what has happened
var Defer = data.MakeProcedure(func(args ...data.Value) (res data.Value) {
	body := args[0].(data.Procedure)
	cleanup := args[1].(data.Procedure)

	defer cleanup.Call()
	return body.Call()
}, 2)

// Read performs the standard LISP read of a string
var Read = data.MakeProcedure(func(args ...data.Value) data.Value {
	v := args[0]
	s := v.(data.Sequence)
	res := read.MustFromString(emptyNamespace, sequence.ToString(s))
	if v, ok := sequence.Last(res); ok {
		return v
	}
	return data.Null
}, 1)
