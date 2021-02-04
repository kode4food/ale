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
	ErrInterfaceCoercionNotSupported = "interface coercion not supported"
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

func (i interfaceWrapper) Wrap(c *Context, v reflect.Value) (data.Value, error) {
	e := v.Elem()
	if !e.IsValid() {
		return data.Nil, nil
	}
	c, err := c.Push(e)
	if err != nil {
		return nil, err
	}
	res := make(data.Object, len(i.methods)+1)
	res[ReceiverKey] = receiver(v)
	for _, m := range i.methods {
		k := data.Keyword(m.name)
		res[k] = m.wrapMethod(v)
	}
	return res, nil
}

func (m *methodWrapper) wrapMethod(v reflect.Value) data.Function {
	switch len(m.out) {
	case 0:
		return m.wrapVoidMethod(v)
	case 1:
		return m.wrapValueMethod(v)
	default:
		return m.wrapVectorMethod(v)
	}
}

func (m *methodWrapper) wrapVoidMethod(v reflect.Value) data.Function {
	inLen := len(m.in)
	fn := v.MethodByName(m.name)

	return data.Applicative(func(in ...data.Value) data.Value {
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			u, err := m.in[i].Unwrap(in[i])
			if err != nil {
				panic(err)
			}
			wIn[i] = u
		}
		fn.Call(wIn)
		return data.Nil
	}, inLen)
}

func (m *methodWrapper) wrapValueMethod(v reflect.Value) data.Function {
	inLen := len(m.in)
	fn := v.MethodByName(m.name)

	return data.Applicative(func(in ...data.Value) data.Value {
		c := &Context{}
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			arg, err := m.in[i].Unwrap(in[i])
			if err != nil {
				panic(err)
			}
			wIn[i] = arg
		}
		wOut := fn.Call(wIn)
		res, err := m.out[0].Wrap(c, wOut[0])
		if err != nil {
			panic(err)
		}
		return res
	}, inLen)
}

func (m *methodWrapper) wrapVectorMethod(v reflect.Value) data.Function {
	inLen := len(m.in)
	outLen := len(m.out)
	fn := v.MethodByName(m.name)

	return data.Applicative(func(in ...data.Value) data.Value {
		c := &Context{}
		wIn := make([]reflect.Value, inLen)
		for i := 0; i < inLen; i++ {
			arg, err := m.in[i].Unwrap(in[i])
			if err != nil {
				panic(err)
			}
			wIn[i] = arg
		}
		wOut := fn.Call(wIn)
		out := make(data.Vector, outLen)
		for i := 0; i < outLen; i++ {
			res, err := m.out[i].Wrap(c, wOut[i])
			if err != nil {
				panic(err)
			}
			out[i] = res
		}
		return out
	}, inLen)
}

func (i interfaceWrapper) Unwrap(v data.Value) (reflect.Value, error) {
	if v, ok := v.(data.Object); ok {
		if r, ok := v[ReceiverKey]; ok {
			if r, ok := r.(receiver); ok {
				res := reflect.Value(r)
				if i.Type != res.Type() {
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
