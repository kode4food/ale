package ffi

import (
	"errors"
	"reflect"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/internal/stream"
)

type channelWrapper struct {
	elem Wrapper
	dir  reflect.ChanDir
}

// ErrChannelCoercionNotSupported is raised when a channel Unwrap is called
var ErrChannelCoercionNotSupported = errors.New(
	"value cannot be coerced into chan",
)

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

func (w *channelWrapper) Wrap(_ *Context, v reflect.Value) (ale.Value, error) {
	r := w.dir&reflect.RecvDir != 0
	s := w.dir&reflect.SendDir != 0
	switch {
	case r && s:
		return data.NewObject(
			data.NewCons(stream.SequenceKey, w.makeSequence(v)),
			data.NewCons(stream.EmitKey, w.makeEmitter(v)),
			data.NewCons(stream.CloseKey, w.makeClose(v)),
		), nil
	case r:
		return data.NewObject(
			data.NewCons(stream.SequenceKey, w.makeSequence(v)),
		), nil
	case s:
		return data.NewObject(
			data.NewCons(stream.EmitKey, w.makeEmitter(v)),
			data.NewCons(stream.CloseKey, w.makeClose(v)),
		), nil
	default:
		return data.EmptyObject, nil
	}
}

func (w *channelWrapper) makeClose(v reflect.Value) data.Procedure {
	return data.MakeProcedure(func(...ale.Value) ale.Value {
		v.Close()
		return data.Null
	}, 0)
}

func (w *channelWrapper) makeSequence(v reflect.Value) data.Sequence {
	var resolver sequence.LazyResolver

	resolver = func() (ale.Value, data.Sequence, bool) {
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
	return data.MakeProcedure(func(args ...ale.Value) ale.Value {
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

func (*channelWrapper) Unwrap(ale.Value) (reflect.Value, error) {
	return _zero, ErrChannelCoercionNotSupported
}
