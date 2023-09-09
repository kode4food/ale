package ffi

import (
	"reflect"

	"github.com/kode4food/ale/data"
)

// Wrappers is a set of Wrapper
type Wrappers []Wrapper

func (w Wrappers) unwrap(args data.Values) []reflect.Value {
	unwrapped := make([]reflect.Value, len(w))
	for i, wrapped := range w {
		arg, err := wrapped.Unwrap(args[i])
		if err != nil {
			panic(err)
		}
		unwrapped[i] = arg
	}
	return unwrapped
}

func (w Wrappers) wrap(args []reflect.Value) data.Values {
	wc := new(Context)
	in := make(data.Values, len(args))
	for i, arg := range args {
		arg, err := w[i].Wrap(wc, arg)
		if err != nil {
			panic(err)
		}
		in[i] = arg
	}
	return in
}
