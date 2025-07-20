package ffi

import (
	"reflect"

	"github.com/kode4food/ale/data"
)

// Wrappers is a set of Wrapper
type Wrappers []Wrapper

func (w Wrappers) mustUnwrap(args data.Vector) []reflect.Value {
	res, err := w.unwrap(args)
	if err != nil {
		panic(err)
	}
	return res
}

func (w Wrappers) unwrap(args data.Vector) ([]reflect.Value, error) {
	unwrapped := make([]reflect.Value, len(w))
	for i, wrapped := range w {
		arg, err := wrapped.Unwrap(args[i])
		if err != nil {
			return nil, err
		}
		unwrapped[i] = arg
	}
	return unwrapped, nil
}

func (w Wrappers) mustWrap(args []reflect.Value) data.Vector {
	res, err := w.wrap(args)
	if err != nil {
		panic(err)
	}
	return res
}

func (w Wrappers) wrap(args []reflect.Value) (data.Vector, error) {
	wc := new(Context)
	in := make(data.Vector, len(args))
	for i, arg := range args {
		arg, err := w[i].Wrap(wc, arg)
		if err != nil {
			return nil, err
		}
		in[i] = arg
	}
	return in, nil
}
