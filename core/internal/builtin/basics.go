package builtin

import (
	"errors"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/read"
)

// Raise will cause a panic
var Raise = data.Applicative(func(args ...data.Value) data.Value {
	err := args[0].(data.String)
	panic(errors.New(string(err)))
}, 1)

// Recover invokes a function and runs a recovery function if Go panics
var Recover = data.Applicative(func(args ...data.Value) (res data.Value) {
	body := args[0].(data.Function)
	rescue := args[1].(data.Function)

	defer func() {
		if rec := recover(); rec != nil {
			err := rec.(error).Error()
			res = rescue.Call(data.String(err))
		}
	}()

	return body.Call()
}, 2)

// Defer invokes a cleanup function, no matter what has happened
var Defer = data.Applicative(func(args ...data.Value) (res data.Value) {
	body := args[0].(data.Function)
	cleanup := args[1].(data.Function)

	defer cleanup.Call()
	return body.Call()
}, 2)

// Read performs the standard LISP read of a string
var Read = data.Applicative(func(args ...data.Value) data.Value {
	v := args[0]
	s := v.(data.Sequence)
	if v, ok := data.Last(read.FromString(sequence.ToStr(s))); ok {
		return v
	}
	return data.Nil
}, 1)
