package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/types"
)

type (
	intfWrapper struct {
		reflect.Type
		methods []*methodWrapper
	}

	receiver reflect.Value

	methodWrapper struct {
		inOutWrappers
		name string
	}
)

const (
	// ErrInterfaceTypeMismatch is raised when an interface of the receiver in
	// a data.Object doesn't match the expected wrapped interface
	ErrInterfaceTypeMismatch = "interface type mismatch"

	// ErrInterfaceCoercionNotSupported is raised when the value to unwrap
	// isn't a data.Object
	ErrInterfaceCoercionNotSupported = "value cannot be coerced into interface"
)

// ReceiverKey is the key used to store an interface receiver
const ReceiverKey = data.Keyword("receiver")

var receiverType = types.MakeBasic("receiver")

func makeWrappedInterface(t reflect.Type) (Wrapper, error) {
	mLen := t.NumMethod()
	res := &intfWrapper{
		Type:    t,
		methods: make([]*methodWrapper, 0, mLen),
	}
	for i := range mLen {
		m := t.Method(i)
		if m.PkgPath != "" { // Not exported
			continue
		}
		w, err := makeWrappedMethod(m)
		if err != nil {
			return nil, err
		}
		res.methods = append(res.methods, w)
	}
	return res, nil
}

func makeWrappedMethod(m reflect.Method) (*methodWrapper, error) {
	res := &methodWrapper{name: m.Name}
	return res, res.wrap(m.Type)
}

func (w *intfWrapper) Wrap(c *Context, v reflect.Value) (ale.Value, error) {
	e := v.Elem()
	if !e.IsValid() {
		return data.Null, nil
	}
	_, err := c.Push(e)
	if err != nil {
		return data.Null, err
	}

	res := make(data.Pairs, len(w.methods)+1)
	res[len(res)-1] = data.NewCons(ReceiverKey, receiver(v))
	for idx, m := range w.methods {
		res[idx] = data.NewCons(
			data.Keyword(m.name),
			m.wrapMethod(v),
		)
	}
	return data.NewObject(res...), nil
}

func (w *methodWrapper) wrapMethod(v reflect.Value) data.Procedure {
	fn := v.MethodByName(w.name)
	return w.wrapFunction(fn)
}

func (w *intfWrapper) Unwrap(v ale.Value) (reflect.Value, error) {
	if res, ok := getReceiver(v); ok {
		if w.Type != res.Type() {
			return _zero, errors.New(ErrInterfaceTypeMismatch)
		}
		return res, nil
	}
	return _zero, errors.New(ErrInterfaceCoercionNotSupported)
}

func getReceiver(v ale.Value) (reflect.Value, bool) {
	if v, ok := v.(*data.Object); ok {
		if r, ok := v.Get(ReceiverKey); ok {
			if r, ok := r.(receiver); ok {
				return reflect.Value(r), true
			}
		}
	}
	return _zero, false
}

func (receiver) Equal(ale.Value) bool {
	return false
}

func (receiver) Type() ale.Type {
	return receiverType
}

func (r receiver) Get(key ale.Value) (ale.Value, bool) {
	return data.DumpMapped(r).Get(key)
}
