package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/data"
)

type (
	interfaceWrapper struct {
		reflect.Type
		methods []*methodWrapper
	}

	receiver reflect.Value

	methodWrapper struct {
		name string
		in   []Wrapper
		out  []Wrapper
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

	// ReceiverType is the type name for an opaque interface receiver
	ReceiverType = data.Name("receiver")
)

func makeWrappedInterface(t reflect.Type) (Wrapper, error) {
	mLen := t.NumMethod()
	res := &interfaceWrapper{
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
	in := make([]Wrapper, cIn)
	for i := 0; i < cIn; i++ {
		w, err := wrapType(t.In(i))
		if err != nil {
			return nil, err
		}
		in[i] = w
	}
	cOut := t.NumOut()
	out := make([]Wrapper, cOut)
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

func (w interfaceWrapper) Wrap(c *Context, v reflect.Value) (data.Value, error) {
	e := v.Elem()
	if !e.IsValid() {
		return data.Nil, nil
	}
	c, err := c.Push(e)
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
	inLen := len(w.in)
	fn := v.MethodByName(w.name)

	return data.Applicative(func(in ...data.Value) data.Value {
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			u, err := w.in[i].Unwrap(in[i])
			if err != nil {
				panic(err)
			}
			wIn[i] = u
		}
		fn.Call(wIn)
		return data.Nil
	}, inLen)
}

func (w *methodWrapper) wrapValueMethod(v reflect.Value) data.Function {
	inLen := len(w.in)
	fn := v.MethodByName(w.name)

	return data.Applicative(func(in ...data.Value) data.Value {
		c := &Context{}
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			arg, err := w.in[i].Unwrap(in[i])
			if err != nil {
				panic(err)
			}
			wIn[i] = arg
		}
		wOut := fn.Call(wIn)
		res, err := w.out[0].Wrap(c, wOut[0])
		if err != nil {
			panic(err)
		}
		return res
	}, inLen)
}

func (w *methodWrapper) wrapVectorMethod(v reflect.Value) data.Function {
	inLen := len(w.in)
	outLen := len(w.out)
	fn := v.MethodByName(w.name)

	return data.Applicative(func(in ...data.Value) data.Value {
		c := &Context{}
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			arg, err := w.in[i].Unwrap(in[i])
			if err != nil {
				panic(err)
			}
			wIn[i] = arg
		}
		wOut := fn.Call(wIn)
		out := make(data.Values, outLen)
		for i := 0; i < outLen; i++ {
			res, err := w.out[i].Wrap(c, wOut[i])
			if err != nil {
				panic(err)
			}
			out[i] = res
		}
		return data.NewVector(out...)
	}, inLen)
}

func (w interfaceWrapper) Unwrap(v data.Value) (reflect.Value, error) {
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

func (receiver) Equal(_ data.Value) bool {
	return false
}

func (receiver) Type() data.Name {
	return ReceiverType
}

func (r receiver) String() string {
	return data.DumpString(r)
}
