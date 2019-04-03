package builtin

import (
	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/read"
	"gitlab.com/kode4food/ale/stdlib"
)

// Raise will cause Go to panic
func Raise(args ...api.Value) api.Value {
	panic(args[0])
}

// Recover invokes a function and runs a recovery function is Go panics
func Recover(args ...api.Value) (res api.Value) {
	body := args[0].(api.Caller).Caller()
	rescue := args[1].(api.Caller).Caller()

	defer func() {
		if rec := recover(); rec != nil {
			res = rescue(rec.(api.Value))
		}
	}()

	return body()
}

// Read performs the standard LISP read of a string
func Read(args ...api.Value) api.Value {
	v := args[0]
	s := v.(api.Sequence)
	if v, ok := api.Last(read.FromString(stdlib.SequenceToStr(s))); ok {
		return v
	}
	return api.Nil
}
