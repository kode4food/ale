package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/types"
)

type (
	intfWrapper struct {
		reflect.Type
		methods []*methodWrapper
	}

	receiver reflect.Value

	methodWrapper struct {
		name string
		in   Wrappers
		out  Wrappers
	}
)

// Error messages
const (
	ErrInterfaceTypeMismatch         = "interface type mismatch"
	ErrInterfaceCoercionNotSupported = "value cannot be coerced into interface"
)

const (
	// ReceiverKey is the key used to store an interface receiver
	ReceiverKey = data.Keyword("receiver")
)

var receiverType = types.Basic("receiver")

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
	t := m.Type
	cIn := t.NumIn()
	in := make(Wrappers, cIn)
	for i := 0; i < cIn; i++ {
		w, err := wrapType(t.In(i))
		if err != nil {
			return nil, err
		}
		in[i] = w
	}
	cOut := t.NumOut()
	out := make(Wrappers, cOut)
	for i := 0; i < cOut; i++ {
		w, err := wrapType(t.Out(i))
		if err != nil {
			return nil, err
		}
		out[i] = w
	}
	return &methodWrapper{
		name: m.Name,
		in:   in,
		out:  out,
	}, nil
}

func (w *intfWrapper) Wrap(c *Context, v reflect.Value) (data.Value, error) {
	e := v.Elem()
	if !e.IsValid() {
		return data.Nil, nil
	}
	_, err := c.Push(e)
	if err != nil {
		return data.Nil, err
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

func (w *methodWrapper) wrapMethod(v reflect.Value) data.Function {
	switch len(w.out) {
	case 0:
		return w.wrapVoidMethod(v)
	case 1:
		return w.wrapValueMethod(v)
	default:
		return w.wrapVectorMethod(v)
	}
}

func (w *methodWrapper) wrapVoidMethod(v reflect.Value) data.Function {
	fn := v.MethodByName(w.name)

	return data.Applicative(func(args ...data.Value) data.Value {
		fn.Call(w.in.unwrap(args))
		return data.Nil
	}, len(w.in))
}

func (w *methodWrapper) wrapValueMethod(v reflect.Value) data.Function {
	fn := v.MethodByName(w.name)

	return data.Applicative(func(args ...data.Value) data.Value {
		in := w.in.unwrap(args)
		out := fn.Call(in)
		res, err := w.out[0].Wrap(new(Context), out[0])
		if err != nil {
			panic(err)
		}
		return res
	}, len(w.in))
}

func (w *methodWrapper) wrapVectorMethod(v reflect.Value) data.Function {
	fn := v.MethodByName(w.name)

	return data.Applicative(func(args ...data.Value) data.Value {
		in := w.in.unwrap(args)
		res := fn.Call(in)
		out := w.out.wrap(res)
		return data.NewVector(out...)
	}, len(w.in))
}

func (w *intfWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if v, ok := v.(data.Object); ok {
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

func (r receiver) String() string {
	return data.DumpString(r)
}
