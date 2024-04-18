package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/internal/stream"
	"github.com/kode4food/ale/pkg/data"
)

type channelWrapper struct {
	elem Wrapper
	dir  reflect.ChanDir
}

// ErrChannelCoercionNotSupported is raised when a channel Unwrap is called
const ErrChannelCoercionNotSupported = "value cannot be coerced into chan"

func makeWrappedChannel(t reflect.Type) (Wrapper, error) {
	w, err := WrapType(t.Elem())
	if err != nil {
		return nil, err
	}
	return &channelWrapper{
		elem: w,
		dir:  t.ChanDir(),
	}, nil
}

func (w *channelWrapper) Wrap(_ *Context, v reflect.Value) (data.Value, error) {
	o := data.NewObject()
	if w.dir&reflect.RecvDir != 0 {
		o = o.Put(data.NewCons(
			stream.SequenceKey, w.makeSequence(v),
		)).(*data.Object)
	}
	if w.dir&reflect.SendDir != 0 {
		o = o.Put(data.NewCons(
			stream.EmitKey, w.makeEmitter(v),
		)).(*data.Object)
		o = o.Put(data.NewCons(
			stream.CloseKey, w.makeClose(v),
		)).(*data.Object)
	}
	return o, nil
}

func (w *channelWrapper) makeClose(v reflect.Value) data.Procedure {
	return data.MakeProcedure(func(...data.Value) data.Value {
		v.Close()
		return data.Null
	}, 0)
}

func (w *channelWrapper) makeSequence(v reflect.Value) data.Sequence {
	var resolver sequence.LazyResolver

	resolver = func() (data.Value, data.Sequence, bool) {
		in, ok := v.Recv()
		if !ok {
			return data.Null, data.EmptyObject, false
		}
		f, err := w.elem.Wrap(new(Context), in)
		if err != nil {
			panic(err)
		}
		return f, sequence.NewLazy(resolver), true
	}

	return sequence.NewLazy(resolver)
}

func (w *channelWrapper) makeEmitter(v reflect.Value) data.Procedure {
	return data.MakeProcedure(func(args ...data.Value) data.Value {
		for _, arg := range args {
			arg, err := w.elem.Unwrap(arg)
			if err != nil {
				panic(err)
			}
			v.Send(arg)
		}
		return data.Null
	})
}

func (*channelWrapper) Unwrap(data.Value) (reflect.Value, error) {
	return _emptyValue, errors.New(ErrChannelCoercionNotSupported)
}
