package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/internal/types"
	"github.com/kode4food/ale/pkg/data"
)

type (
	intfWrapper struct {
		reflect.Type
		methods []*methodWrapper
	}

	receiver reflect.Value

	methodWrapper struct {
		*inOutWrappers
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
	for i := 0; i < mLen; i++ {
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
	io, err := makeInOutWrappers(m.Type)
	if err != nil {
		return nil, err
	}
	return &methodWrapper{
		name:          m.Name,
		inOutWrappers: io,
	}, nil
}

func (w *intfWrapper) Wrap(c *Context, v reflect.Value) (data.Value, error) {
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

func (w *intfWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if v, ok := v.(*data.Object); ok {
		if r, ok := v.Get(ReceiverKey); ok {
			if r, ok := r.(receiver); ok {
				res := reflect.Value(r)
				if w.Type != res.Type() {
					return _emptyValue, errors.New(ErrInterfaceTypeMismatch)
				}
				return res, nil
			}
		}
	}
	return _emptyValue, errors.New(ErrInterfaceCoercionNotSupported)
}

func (receiver) Equal(data.Value) bool {
	return false
}

func (receiver) Type() types.Type {
	return receiverType
}

func (r receiver) Get(key data.Value) (data.Value, bool) {
	return data.DumpMapped(r).Get(key)
}
