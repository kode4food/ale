package builtin

import (
	"errors"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/read"
	"gitlab.com/kode4food/ale/stdlib"
)

// Raise will cause a panic
func Raise(args ...data.Value) data.Value {
	err := args[0].(data.String)
	panic(errors.New(string(err)))
}

// Recover invokes a function and runs a recovery function if Go panics
func Recover(args ...data.Value) (res data.Value) {
	body := args[0].(data.Caller).Caller()
	rescue := args[1].(data.Caller).Caller()

	defer func() {
		if rec := recover(); rec != nil {
			err := rec.(error).Error()
			res = rescue(data.String(err))
		}
	}()

	return body()
}

// Defer invokes a cleanup function, no matter what has happened
func Defer(args ...data.Value) (res data.Value) {
	body := args[0].(data.Caller).Caller()
	cleanup := args[1].(data.Caller).Caller()

	defer cleanup()
	return body()
}

// Read performs the standard LISP read of a string
func Read(args ...data.Value) data.Value {
	v := args[0]
	s := v.(data.Sequence)
	if v, ok := data.Last(read.FromString(stdlib.SequenceToStr(s))); ok {
		return v
	}
	return data.Nil
}